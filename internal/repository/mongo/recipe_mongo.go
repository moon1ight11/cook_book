package mongo

import (
    "cook_book/internal/repository"
    "log/slog"
    "go.mongodb.org/mongo-driver/mongo"
)

type recipeMongoRepo struct {
    collection *mongo.Collection
    logger     *slog.Logger
}

func NewRecipeMongoRepo(db *mongo.Database, logger *slog.Logger) repository.RecipeRepository {
    return &recipeMongoRepo{
        collection: db.Collection("recipes"),
        logger:     logger,
    }
}