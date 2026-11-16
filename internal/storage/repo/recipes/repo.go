package recipes

import (
	"cook_book/internal/storage"
	"cook_book/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type recipeRepo struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewRecipeMongoRepo(db *storage.DataBase, logger logger.Logger) *recipeRepo {
	return &recipeRepo{
		collection: db.DB.Collection("recipes"),
		logger:     logger,
	}
}
