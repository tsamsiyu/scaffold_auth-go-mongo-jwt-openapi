package schema

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddUsersIndexes(ctx context.Context, db *mongo.Database) error {
	if _, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true).SetName("uniq_email"),
	}); err != nil {
		return err
	}

	return nil
}

func Migrate(ctx context.Context, db *mongo.Database) error {
	if err := AddUsersIndexes(ctx, db); err != nil {
		return err
	}

	return nil
}
