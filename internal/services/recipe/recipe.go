package recipe

import (
	"cook_book/internal/storage/cache"
	"cook_book/internal/storage/repo"
)

// структура юзкейсов
type recipeUseCase struct {
	repo  repo.RecipeRepositoryInterface
	cache cache.RecipeCacheInterface
}

// коструктор
func NewRecipeUseCase(repo repo.RecipeRepositoryInterface, cache cache.RecipeCacheInterface) *recipeUseCase {
	return &recipeUseCase{
		repo:  repo,
		cache: cache,
	}
}
