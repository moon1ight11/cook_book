package recipe

import (
    "context"
    "cook_book/internal/cache"
    "cook_book/internal/domain"
    "cook_book/internal/repository"
)

// интерфейс для юзкейсов рецептов
type RecipeUseCase interface {
    CreateRecipe(ctx context.Context, recipe *domain.Recipe) error
    UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error
    DeleteRecipe(ctx context.Context, id string) error
    GetRecipeByID(ctx context.Context, id string) (*domain.Recipe, error)
    ListRecipes(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domain.Recipe, int64, error)
    SearchByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error)
    GetTotalCount(ctx context.Context) (int64, error)
}

// структура юзкейсов
type recipeUseCase struct {
    repo  repository.RecipeRepository
    cache cache.RecipeCache
}

// коструктор
func NewRecipeUseCase(repo repository.RecipeRepository, cache cache.RecipeCache) RecipeUseCase {
    return &recipeUseCase{
        repo:  repo,
        cache: cache,
    }
}