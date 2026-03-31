package storage

import "go.mongodb.org/mongo-driver/mongo"

type DataBase struct {
	DB *mongo.Database
}
