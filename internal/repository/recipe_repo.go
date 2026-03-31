package repository

import (
    "context"
    "cook_book/internal/domain"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeRepository interface {
    Create(ctx context.Context, recipe *domain.Recipe) error
    Update(ctx context.Context, recipe *domain.Recipe) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Recipe, error)
    FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int64) ([]*domain.Recipe, int64, error)
    FindByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error)
}