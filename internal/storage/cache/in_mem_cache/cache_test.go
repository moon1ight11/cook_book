package inmemcache

import (
	"cook_book/internal/domain"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// функция для создания тестового кэша
func createTestCache(ttl time.Duration) *inMemoryCache {
	return &inMemoryCache{
		recipes:         make(map[string]*cacheItem),
		tagRecipes:      make(map[string]*tagCacheItem),
		ttl:             ttl,
		cleanupInterval: time.Minute * 15,
		stopCleanup:     make(chan bool),
		mu:              sync.RWMutex{},
		logger:          slog.Default(),
	}
}

// функция для создания тестового рецепта
func createTestRecipe(title string) *domain.Recipe {
	return &domain.Recipe{
		Id:          primitive.NewObjectID(),
		Title:       title,
		Description: "Test description",
		Difficulty:  "medium",
		Tags:        []string{"test", "simple"},
	}
}

// установка и получение кэша
func TestCache_SetAndGet(t *testing.T) {
	cache := createTestCache(time.Minute)

	recipe := createTestRecipe("TestRecipe1")
	id := recipe.Id.Hex()

	// устанавливаем в кэш
	cache.Set(id, recipe)

	// получаем из кэша
	cached, found := cache.Get(id)

	// проверки
	require.True(t, found, "Recipe must be found")
	assert.Equal(t, recipe.Id, cached.Id)
	assert.Equal(t, recipe.Title, cached.Title)
	assert.Equal(t, recipe.Description, cached.Description)
	assert.Equal(t, recipe.Tags, cached.Tags)
}

// промах в кэше
func TestCache_Get_Miss(t *testing.T) {
	cache := createTestCache(time.Minute)

	// пытаемся получить несуществующий рецепт
	cached, found := cache.Get("FAKEID")

	assert.False(t, found)
	assert.Nil(t, cached)
}

// получение просроченного кэша
func TestCache_Get_Expired(t *testing.T) {
	// создаем кэш с коротким ttl
	cache := createTestCache(time.Millisecond)
	
	recipe := createTestRecipe("TestRecipe2")
	id := recipe.Id.Hex()

	// устанавливаем в кэш
	cache.Set(id, recipe)

	// ждем, пока истечет
	time.Sleep(10 * time.Millisecond)

	// пытаемся получить просроченный рецепт
	cached, found := cache.Get(id)

	assert.False(t, found, "Recipe cant be found")
	assert.Nil(t, cached)
}

// удаление записи из кэша
func TestCache_Delete(t *testing.T) {
	cache := createTestCache(time.Minute)
	
	recipe := createTestRecipe("TestRecipe3")
	id := recipe.Id.Hex()

	// добавляем и проверяем, что есть
	cache.Set(id, recipe)
	_, found := cache.Get(id)
	require.True(t, found, "Recipe must be found")

	// удаляем
	cache.Delete(id)

	// проверяем, что удалился
	_, found = cache.Get(id)
	assert.False(t, found, "Recipe cant be found")
}

// удаление несуществующего кэша
func TestCache_Delete_NonExistent(t *testing.T) {
	cache := createTestCache(time.Minute)

	// удаляем несуществующий кэш
	cache.Delete("FAKEID")
	
	// проверяем, что все работает
	_, found := cache.Get("FAKEID")
	assert.False(t, found)
}

// очистка
func TestCache_Clear(t *testing.T) {
	cache := createTestCache(time.Minute)

	// добавляем несколько записей
	for i := 0; i < 5; i++ {
		recipe := createTestRecipe("Recipe")
		recipe.Id = primitive.NewObjectID()
		cache.Set(recipe.Id.Hex(), recipe)
	}

	// добавляем записи с тегами
	tagsKey := "test,simple"
	recipes := []*domain.Recipe{
		createTestRecipe("TaggedRecipe1"),
		createTestRecipe("TaggedRecipe2"),
	}
	cache.SetByTags(tagsKey, recipes)

	// проверяем что все добавилось
	assert.Len(t, cache.recipes, 5)
	assert.Len(t, cache.tagRecipes, 1)

	// очищаем кэш
	cache.Clear()

	// проверяем что все очистилось
	assert.Len(t, cache.recipes, 0)
	assert.Len(t, cache.tagRecipes, 0)
}

// установка и получение кэша по тегам
func TestCache_SetByTagsAndGetByTags(t *testing.T) {
	cache := createTestCache(time.Minute)

	tagsKey := "tag1,tag2"
	recipes := []*domain.Recipe{
		createTestRecipe("Recipe5"),
		createTestRecipe("Recipe6"),
		createTestRecipe("Recipe7"),
	}

	// устанавливаем по тегам
	cache.SetByTags(tagsKey, recipes)

	// получаем по тегам
	cached, found := cache.GetByTags(tagsKey)

	require.True(t, found)
	assert.Len(t, cached, 3)
	assert.Equal(t, recipes[0].Title, cached[0].Title)
	assert.Equal(t, recipes[1].Title, cached[1].Title)
	assert.Equal(t, recipes[2].Title, cached[2].Title)
}

// промах по тегам
func TestCache_GetByTags_Miss(t *testing.T) {
	cache := createTestCache(time.Minute)

	// пытаемся получить несуществующие теги
	cached, found := cache.GetByTags("TAGNOTFOUND")

	assert.False(t, found)
	assert.Nil(t, cached)
}

// получение просроченного кэша по тегам
func TestCache_GetByTags_Expired(t *testing.T) {
	cache := createTestCache(time.Millisecond)

	tagsKey := "tag4,tag5"
	recipes := []*domain.Recipe{
		createTestRecipe("Recipe8"),
	}

	cache.SetByTags(tagsKey, recipes)

	// ждем пока истечет
	time.Sleep(10 * time.Millisecond)

	// пытаемся получить просроченное
	cached, found := cache.GetByTags(tagsKey)

	assert.False(t, found)
	assert.Nil(t, cached)
}

