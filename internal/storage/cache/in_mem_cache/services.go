package inmemcache

import (
	"time"

	"cook_book/internal/domain"
)

func (c *inMemoryCache) Get(id string) (*domain.Recipe, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.recipes[id]
	if !exists {
		c.logger.Info("Cache miss", "id", id)
		return nil, false
	}
	c.logger.Info("Cache hit", "id", id)

	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.recipe, true
}

func (c *inMemoryCache) Set(id string, recipe *domain.Recipe) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.recipes[id] = &cacheItem{
		recipe:    recipe,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *inMemoryCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.recipes, id)
}

func (c *inMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.recipes = make(map[string]*cacheItem)
	c.tagRecipes = make(map[string]*tagCacheItem)
}

func (c *inMemoryCache) GetByTags(tagsKey string) ([]*domain.Recipe, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.tagRecipes[tagsKey]
	if !exists {
		c.logger.Info("Cache by tags miss", "tags", tagsKey)
		return nil, false
	}
	c.logger.Info("Cache by tags hit", "tags", tagsKey)

	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.recipes, true
}

func (c *inMemoryCache) SetByTags(tagsKey string, recipes []*domain.Recipe) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tagRecipes[tagsKey] = &tagCacheItem{
		recipes:   recipes,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *inMemoryCache) cleanupExpired() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()

			for key, item := range c.recipes {
				if now.After(item.expiresAt) {
					delete(c.recipes, key)
				}
			}

			for key, item := range c.tagRecipes {
				if now.After(item.expiresAt) {
					delete(c.tagRecipes, key)
				}
			}

			c.mu.Unlock()

		case <-c.stopCleanup:
			return
		}
	}
}

func (c *inMemoryCache) Stop() {
	close(c.stopCleanup)
	c.wg.Wait()
}
