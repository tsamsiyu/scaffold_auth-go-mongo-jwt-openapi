package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStatus string

const (
	collectionName = "users"
)

var (
	Pending UserStatus = "pending"
	Active  UserStatus = "active"
)

type SignUpRequest struct {
	Token      string     `bson:"token"`
	Code       string     `bson:"code"`
	CreatedAt  time.Time  `bson:"createdAt"`
	NotifiedAt *time.Time `bson:"notifiedAt"`
}

type User struct {
	UID          string         `bson:"uid"`
	Name         string         `bson:"name"`
	Email        string         `bson:"email"`
	Status       UserStatus     `bson:"status"`
	PasswordHash string         `bson:"passwordHash"`
	LastLoggedIn *time.Time     `bson:"lastLoggedIn"`
	CreatedAt    time.Time      `bson:"createdAt"`
	ConfirmedAt  *time.Time     `bson:"confirmedAt"`
	SignUpReq    *SignUpRequest `bson:"signUpReq"`
}

type UserRepository interface {
	FindAllNotNotifiedSignUpRequests(ctx context.Context) ([]User, error)
	SaveNotifiedSignUpReqTime(ctx context.Context, uid string, t time.Time) error
	FindBySignUpReqToken(ctx context.Context, token string) (*User, error)
	Create(ctx context.Context, model *User) error
	ConfirmAndDeleteSignUpReq(ctx context.Context, uid string) error
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
	cursor, err := r.db.Collection(collectionName).Find(ctx, bson.M{
		"signUpReq.notifiedAt": nil,
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var models []User

	for cursor.Next(ctx) {
		var model User

		err := cursor.Decode(&model)
		if err != nil {
			return nil, err
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (r *mongoUserRepository) SaveNotifiedSignUpReqTime(ctx context.Context, uid string, t time.Time) error {
	_, err := r.db.Collection(collectionName).UpdateOne(ctx, bson.M{
		"uid": uid,
	}, bson.M{
		"$set": bson.M{"signUpReq.notifiedAt": t.Format(time.RFC3339)},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoUserRepository) FindBySignUpReqToken(ctx context.Context, token string) (*User, error) {
	singleResult := r.db.Collection(collectionName).FindOne(ctx, bson.M{
		"signUpReq.token": token,
	})
	if err := singleResult.Err(); err != nil {
		return nil, err
	}

	var u User

	if err := singleResult.Decode(&u); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mongoUserRepository) ConfirmAndDeleteSignUpReq(ctx context.Context, uid string) error {
	_, err := r.db.Collection(collectionName).UpdateOne(ctx, bson.M{
		"uid": uid,
	}, bson.M{
		"$set": bson.M{"signUpReq": nil, "status": Active},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoUserRepository) FindByUID(ctx context.Context, uid string) (*User, error) {
	singleResult := r.db.Collection(collectionName).FindOne(ctx, bson.M{
		"uid": uid,
	})
	if err := singleResult.Err(); err != nil {
		return nil, err
	}

	var u User

	if err := singleResult.Decode(&u); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	singleResult := r.db.Collection(collectionName).FindOne(ctx, bson.M{
		"email": email,
	})
	if err := singleResult.Err(); err != nil {
		return nil, err
	}

	var u User

	if err := singleResult.Decode(&u); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *mongoUserRepository) Create(ctx context.Context, model *User) error {
	doc, err := bson.Marshal(model)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(collectionName).InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}
