package recipes

import (
	"cook_book/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type recipeRepo struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewRecipeMongoRepo(db *mongo.Database, logger logger.Logger) *recipeRepo {
	return &recipeRepo{
		collection: db.Collection("recipes"),
		logger:     logger,
	}
}
