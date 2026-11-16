package cache

import (
	"cook_book/internal/domain"
)

type RecipeCacheInterface interface {
	Get(id string) (*domain.Recipe, bool)
	Set(id string, recipe *domain.Recipe)
	Delete(id string)
	Clear()
	GetByTags(tagsKey string) ([]*domain.Recipe, bool)
	SetByTags(tagsKey string, recipes []*domain.Recipe)
	Stop()
}

type CacheStats struct {
	TotalItems   int
	TotalTagSets int
}
