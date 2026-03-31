package cache

import (
	"cook_book/internal/domain"
)

// интерфейс для кэша
type RecipeCache interface {
	Get(id string) (*domain.Recipe, bool)
	Set(id string, recipe *domain.Recipe)
	Delete(id string)
	Clear()
	GetByTags(tagsKey string) ([]*domain.Recipe, bool)
	SetByTags(tagsKey string, recipes []*domain.Recipe)
	GetStats() CacheStats
}

// статистика кэша
type CacheStats struct {
    TotalItems   int
    TotalTagSets int
}