package storage

import (
	"context"
	"cook_book/internal/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewStorage(cfg *config.Config) (*DataBase, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin",
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
	)
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error in NewStorage: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error in NewStorage: %w", err)
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	Database := DataBase{
		DB: db,
		Client: mongoClient,
	}

	return &Database, nil
}

