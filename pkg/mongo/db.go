package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MainDB *mongo.Database

func NewClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return client, err
}

func ProvideDatabase(client *mongo.Client, dbname string) MainDB {
	return client.Database(dbname)
}
