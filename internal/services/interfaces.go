package services

import (
	"context"
	"cook_book/internal/domain"
)

// интерфейс для юзкейсов рецептов
type RecipeUseCaseInterface interface {
	CreateRecipe(ctx context.Context, recipe *domain.Recipe) error
	UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error
	DeleteRecipe(ctx context.Context, id string) error
	GetRecipeByID(ctx context.Context, id string) (*domain.Recipe, error)
	ListRecipes(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domain.Recipe, int64, error)
	SearchByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error)
	GetTotalCount(ctx context.Context) (int64, error)
}
