package recipe

import (
	"cook_book/internal/storage/cache"
	"cook_book/internal/storage/repo"
)

type recipeUseCase struct {
	repo  repo.RecipeRepositoryInterface
	cache cache.RecipeCacheInterface
}

func NewRecipeUseCase(repo repo.RecipeRepositoryInterface, cache cache.RecipeCacheInterface) *recipeUseCase {
	return &recipeUseCase{
		repo:  repo,
		cache: cache,
	}
}
