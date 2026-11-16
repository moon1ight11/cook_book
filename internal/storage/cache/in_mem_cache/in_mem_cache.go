package inmemcache

import (
	"cook_book/internal/domain"
	"cook_book/internal/storage/cache"
	"cook_book/pkg/logger"
	"sync"
	"time"
)

type inMemoryCache struct {
	recipes         map[string]*cacheItem
	tagRecipes      map[string]*tagCacheItem
	logger          logger.Logger
	mu              sync.RWMutex
	wg              sync.WaitGroup
	ttl             time.Duration
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

type cacheItem struct {
	recipe    *domain.Recipe
	expiresAt time.Time
}

type tagCacheItem struct {
	recipes   []*domain.Recipe
	expiresAt time.Time
}

func NewRecipeCache(ttl, cleanupInterval time.Duration, logger logger.Logger) cache.RecipeCacheInterface {
	cache := &inMemoryCache{
		recipes:         make(map[string]*cacheItem),
		tagRecipes:      make(map[string]*tagCacheItem),
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan bool),
		logger:          logger,
	}

	cache.wg.Add(1)
	go func() {
		defer cache.wg.Done()
		cache.cleanupExpired()
	}()

	return cache
}
