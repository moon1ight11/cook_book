package inmemcache

import (
	"cook_book/internal/storage/cache"
	"cook_book/internal/domain"
	"sync"
	"time"
)

// структура кэша
type inMemoryCache struct {
	recipes         map[string]*cacheItem
	tagRecipes      map[string]*tagCacheItem
	mu              sync.RWMutex
	ttl             time.Duration
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

// структура единицы хранения кэша
type cacheItem struct {
	recipe    *domain.Recipe
	expiresAt time.Time
}

// структура единицы хранения кэша нескольких элементов по тегу
type tagCacheItem struct {
	recipes   []*domain.Recipe
	expiresAt time.Time
}

// конструктор кэша
func NewRecipeCache(ttl, cleanupInterval time.Duration) cache.RecipeCacheInterface {
	cache := &inMemoryCache{
		recipes:         make(map[string]*cacheItem),
		tagRecipes:      make(map[string]*tagCacheItem),
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan bool),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}
