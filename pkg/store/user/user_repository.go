package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserStatus string

var (
	Pending UserStatus = "pending"
	Active  UserStatus = "active"
)

type User struct {
	UID          string
	Name         string
	Email        string
	Status       UserStatus
	PasswordHash string
	LastLoggedIn *time.Time
	Created      time.Time

	SignUpConfirmationCode *string
}

type UserRepository interface {
	GetByUID(ctx context.Context, uid string) (*User, error)
	Create(ctx context.Context, model *User) error
}

type mongoUserRepository struct {
	client *mongo.Client
}

func (r *mongoUserRepository) GetByUID(ctx context.Context, uid string) (*User, error) {
	return nil, nil
}

func (r *mongoUserRepository) Create(ctx context.Context, model *User) error {
	return nil
}
