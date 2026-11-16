package users

import (
	"cook_book/internal/storage"
	"cook_book/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewUserMongoRepo(db *storage.DataBase, logger logger.Logger) *userRepo {
	return &userRepo{
		collection: db.DB.Collection("users"),
		logger:     logger,
	}
}
