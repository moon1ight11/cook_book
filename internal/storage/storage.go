package storage

import (
	"context"
	"cook_book/internal/config"
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStorage(cfg *config.Config) (*mongo.Database, error) {
	mongoURI := fmt.Sprintf("")
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	return db, nil
}
