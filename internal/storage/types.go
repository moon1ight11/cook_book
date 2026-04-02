package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type DataBase struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func (db *DataBase) Close(ctx context.Context) error {
	if db.Client != nil {
		return db.Client.Disconnect(ctx)
	}
	return nil
}
