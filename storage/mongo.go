package storage

import (
	"context"
	"fmt"
	"time"

	// "my-app/types"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoStore struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewDatabase(ctx context.Context, uri string, database string, collection string) (*mongoStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDb client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(database)
	collectionName := db.Collection(collection)

	return &mongoStore{
		client:     client,
		db:         db,
		collection: collectionName,
	}, nil
}

func (db *mongoStore) CloseStore(ctx context.Context) error {
	if db.client != nil {
		err := db.client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("failed to close MongoDB connection: %v", err)
		}
	}

	return nil
}
