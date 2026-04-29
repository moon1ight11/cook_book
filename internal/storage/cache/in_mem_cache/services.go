package inmemcache

import (
	"cook_book/internal/domain"
	"time"
)

// получение записи по id
func (c *inMemoryCache) Get(id string) (*domain.Recipe, bool) {
	// блокируем чтение
	c.mu.RLock()
	defer c.mu.RUnlock()

	// получаем рецепт
	item, exists := c.recipes[id]
	if !exists {
		c.logger.Info("Cache miss", "id", id)
		return nil, false
	}
	c.logger.Info("Cache hit", "id", id)

	// если запись просрочена - не возвращаем ничего
	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.recipe, true
}

// установка записи в кэш
func (c *inMemoryCache) Set(id string, recipe *domain.Recipe) {
	// блокируем запись
	c.mu.Lock()
	defer c.mu.Unlock()

	// устанавливаем кэш
	c.recipes[id] = &cacheItem{
		recipe:    recipe,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// удаление записи
func (c *inMemoryCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.recipes, id)
}

// очистка кэша
func (c *inMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.recipes = make(map[string]*cacheItem)
	c.tagRecipes = make(map[string]*tagCacheItem)
}

// получение записей по тегам
func (c *inMemoryCache) GetByTags(tagsKey string) ([]*domain.Recipe, bool) {
	// блокируем чтение
	c.mu.RLock()
	defer c.mu.RUnlock()

	// получаем записи
	item, exists := c.tagRecipes[tagsKey]
	if !exists {
		c.logger.Info("Cache by tags miss", "tags", tagsKey)
		return nil, false
	}
	c.logger.Info("Cache by tags hit", "tags", tagsKey)

	// если просрочено - прокидываем
	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.recipes, true
}

// установка по тегам
func (c *inMemoryCache) SetByTags(tagsKey string, recipes []*domain.Recipe) {
	// блокируем запись
	c.mu.Lock()
	defer c.mu.Unlock()

	// устанавливаем
	c.tagRecipes[tagsKey] = &tagCacheItem{
		recipes:   recipes,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// периодическая очистка просроченного кэша
func (c *inMemoryCache) cleanupExpired() {
	// создаем тикер с заданным интервалом
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		// если срабатывает тикер
		case <-ticker.C:
			// блокируем запись
			c.mu.Lock()
			now := time.Now()

			// очищаем единичные записи, если просрочены
			for key, item := range c.recipes {
				if now.After(item.expiresAt) {
					delete(c.recipes, key)
				}
			}

			// очищаем массивы записей с тегами
			for key, item := range c.tagRecipes {
				if now.After(item.expiresAt) {
					delete(c.tagRecipes, key)
				}
			}

			c.mu.Unlock()

		// если срабатывает стопОчистка
		case <-c.stopCleanup:
			return
		}
	}
}

// остановка очистки кэша
func (c *inMemoryCache) Stop() {
	close(c.stopCleanup)
	c.wg.Wait()
}
