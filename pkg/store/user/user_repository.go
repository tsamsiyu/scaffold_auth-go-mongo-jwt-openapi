package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDuplicateError struct {
	error
}

type UserStatus string

const (
	CollectionName = "users"
)

var (
	StatusPending   UserStatus = "pending"
	StatusConfirmed UserStatus = "confirmed"
)

type SignUpRequest struct {
	Token      string     `bson:"token"`
	Code       string     `bson:"code"`
	NotifiedAt *time.Time `bson:"notifiedAt"`
}

type User struct {
	UID          string         `bson:"_id"`
	Name         string         `bson:"name"`
	Email        string         `bson:"email"`
	Status       UserStatus     `bson:"status"`
	PasswordHash string         `bson:"passwordHash"`
	CreatedAt    time.Time      `bson:"createdAt"`
	ConfirmedAt  *time.Time     `bson:"confirmedAt"`
	SignUpReq    *SignUpRequest `bson:"signUpReq"`
}

type UserRepository interface {
	FindAllNotNotifiedSignUpRequests(ctx context.Context) ([]User, error)
	DeleteAllPendingOlderThan(ctx context.Context, t time.Time) (int, error)
	SaveNotifiedSignUpReqTime(ctx context.Context, uid string, t time.Time) error
	FindBySignUpReqToken(ctx context.Context, token string) (*User, error)
	Create(ctx context.Context, model *User) error
	ConfirmAndDeleteSignUpReq(ctx context.Context, uid string) (bool, error)
	FindByUID(ctx context.Context, uid string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type mongoUserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{
		db: db,
	}
}

func (r *mongoUserRepository) FindAllNotNotifiedSignUpRequests(ctx context.Context) ([]User, error) {
	cursor, err := r.db.Collection(CollectionName).Find(ctx, bson.D{
		{"signUpReq.notifiedAt", nil},
		{"signUpReq", bson.M{"$ne": nil}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	models := make([]User, 0)

	for cursor.Next(ctx) {
		var model User

		err := cursor.Decode(&model)
		if err != nil {
			return nil, err
		}

		models = append(models, model)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (r *mongoUserRepository) DeleteAllPendingOlderThan(ctx context.Context, t time.Time) (int, error) {
	res, err := r.db.Collection(CollectionName).DeleteMany(ctx, bson.D{
		{"status", StatusPending},
		{"createdAt", bson.M{"$lt": t}},
	})
	if err != nil {
		return 0, err
	}

	return int(res.DeletedCount), nil
}

func (r *mongoUserRepository) SaveNotifiedSignUpReqTime(ctx context.Context, uid string, t time.Time) error {
	_, err := r.db.Collection(CollectionName).UpdateOne(ctx, bson.M{
		"_id": uid,
	}, bson.M{
		"$set": bson.M{"signUpReq.notifiedAt": t},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoUserRepository) FindBySignUpReqToken(ctx context.Context, token string) (*User, error) {
	singleResult := r.db.Collection(CollectionName).FindOne(ctx, bson.D{
		{"signUpReq.token", token},
	})
	if err := singleResult.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	var u User

	if err := singleResult.Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mongoUserRepository) ConfirmAndDeleteSignUpReq(ctx context.Context, uid string) (bool, error) {
	res, err := r.db.Collection(CollectionName).UpdateOne(ctx, bson.D{
		{"_id", uid},
		{"status", StatusPending},
	}, bson.M{
		"$set": bson.M{
			"signUpReq":   nil,
			"confirmedAt": time.Now(),
			"status":      StatusConfirmed,
		},
	})
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

func (r *mongoUserRepository) FindByUID(ctx context.Context, uid string) (*User, error) {
	singleResult := r.db.Collection(CollectionName).FindOne(ctx, bson.D{
		{"_id", uid},
	})
	if err := singleResult.Err(); err != nil {
		return nil, err
	}

	var u User

	if err := singleResult.Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	singleResult := r.db.Collection(CollectionName).FindOne(ctx, bson.D{
		{"email", email},
	})
	if err := singleResult.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	var u User

	if err := singleResult.Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mongoUserRepository) Create(ctx context.Context, model *User) error {
	doc, err := bson.Marshal(model)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(CollectionName).InsertOne(ctx, doc)
	if err != nil {
		return mapError(err)
	}

	return nil
}

func mapError(err error) error {
	if mongo.IsDuplicateKeyError(err) {
		return &UserDuplicateError{}
	}

	return err
}
