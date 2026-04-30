package recipes

import (
	"context"
	"cook_book/internal/domain"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testRepo *recipeRepo
var testClient *mongo.Client

// точка входа в тестирование
func TestMain(m *testing.M) {
	ctx := context.Background()

	// соединяемся с тестовой бд
	mongoURI := "mongodb://admin:123@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	// создаем тестовые бд и коллекцию пользователей
	testClient = client
	testDB := client.Database("cook_book_test")
	testCollection := testDB.Collection("recipes_test")

	// создаем репозиторий
	testRepo = &recipeRepo{
		collection: testCollection,
		logger:     slog.Default(),
	}

	// запускаем тесты
	code := m.Run()

	// после тестов все удаляем
	testCollection.Drop(ctx)
	client.Disconnect(ctx)

	// выходим из тестов
	os.Exit(code)
}

// хэлпер для создания тестового рецепта
func createTestRecipe(t *testing.T, title string) *domain.Recipe {
	ctx := context.Background()

	recipe := &domain.Recipe{
		Title:       title,
		Description: "Test description",
		OwnerId:     primitive.NewObjectID(),
		Ingredients: []domain.Ingredient{
			{Name: "Ingredient 1", Quantity: "100g"},
			{Name: "Ingredient 2", Quantity: "200ml"},
		},
		Steps:      []string{"Step 1", "Step 2", "Step 3"},
		PrepTime:   15,
		CookTime:   30,
		Servings:   4,
		Difficulty: "medium",
		Tags:       []string{"test", "simple"},
	}

	err := testRepo.Create(ctx, recipe)
	require.NoError(t, err)

	return recipe
}

// тест на создание рецепта
func TestRecipeRepo_Create(t *testing.T) {
	ctx := context.Background()

	// задаем рецепт
	recipe := &domain.Recipe{
		Title:       "TestRecipe",
		Description: "Test Description",
		OwnerId:     primitive.NewObjectID(),
		Ingredients: []domain.Ingredient{
			{Name: "One", Quantity: "200g"},
			{Name: "Two", Quantity: "100g"},
		},
		Steps:      []string{"First", "Second"},
		PrepTime:   10,
		CookTime:   30,
		Servings:   4,
		Difficulty: "easy",
		Tags:       []string{"dessert", "baking"},
	}

	// добавляем в репо
	err := testRepo.Create(ctx, recipe)

	// проверка выходных данных
	assert.NoError(t, err)
	assert.NotEmpty(t, recipe.Id)
	assert.False(t, recipe.CreatedAt.IsZero())
}

// тест на поиск рецепта по id
func TestRecipeRepo_FindByID(t *testing.T) {
	ctx := context.Background()

	// создаем тестовый рецепт
	recipe := createTestRecipe(t, "TestRecipe1")

	// ищем по id
	found, err := testRepo.FindByID(ctx, recipe.Id)

	// проверка выходных данных
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, recipe.Title, found.Title)
	assert.Equal(t, recipe.Description, found.Description)
	assert.Equal(t, recipe.OwnerId, found.OwnerId)

}

// тест на поиск по айди - не найдено
func TestRecipeRepo_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()

	// делаем фейковый айди
	fakeID := primitive.NewObjectID()

	// ищем по нему
	_, err := testRepo.FindByID(ctx, fakeID)

	// не должны ничего найти
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// тест на поиск всех рецептов с фильтрами и пагинацией
func TestRecipeRepo_FindAll(t *testing.T) {
	ctx := context.Background()

	// создаем несколько рецептов
	createTestRecipe(t, "TestRecipe3")
	createTestRecipe(t, "TestRecipe4")
	createTestRecipe(t, "TestRecipe5")

	// просто получаем все что есть
	recipes, total, err := testRepo.FindAll(ctx, map[string]interface{}{}, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, recipes)
	assert.GreaterOrEqual(t, total, int64(3))

	// с фильтром по названию
	filter := map[string]interface{}{
		"title": "TestRecipe3",
	}
	recipes, total, err = testRepo.FindAll(ctx, filter, 10, 0)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	if len(recipes) > 0 {
		assert.Contains(t, recipes[0].Title, "TestRecipe3")
	}

	// с пагинацией
	recipes, total, err = testRepo.FindAll(ctx, map[string]interface{}{}, 2, 0)
	assert.NoError(t, err)
	assert.LessOrEqual(t, int64(len(recipes)), int64(2))
	assert.GreaterOrEqual(t, total, int64(3))
}

// тест на поиск всех рецептов с фильтром на сложность
func TestRecipeRepo_FindAll_WithDifficultyFilter(t *testing.T) {
	ctx := context.Background()

	// создаем рецепты с разной сложностью
	createTestRecipe(t, "TestRecipe6")

	filter := map[string]interface{}{
		"difficulty": "medium",
	}

	recipes, total, err := testRepo.FindAll(ctx, filter, 10, 0)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, recipe := range recipes {
		assert.Equal(t, "medium", recipe.Difficulty)
	}
}

// тест на поиск всех рецептов по тегу
func TestRecipeRepo_FindByTags(t *testing.T) {
	ctx := context.Background()

	// создаем рецепты
	createTestRecipe(t, "TestRecipe7")
	createTestRecipe(t, "TestRecipe8")

	// ищем по тегу тест
	recipes, err := testRepo.FindByTags(ctx, []string{"test"})

	// сверяем че нашли
	assert.NoError(t, err)
	assert.NotNil(t, recipes)
	assert.GreaterOrEqual(t, len(recipes), 2)
}

// тест на поиск всех рецептов по несуществующему тегу
func TestRecipeRepo_FindByTags_NotFound(t *testing.T) {
	ctx := context.Background()

	// ищем по несуществующему тегу
	recipes, err := testRepo.FindByTags(ctx, []string{"TAGNOTFOUND"})

	// не должны найти
	assert.NoError(t, err)
	assert.Empty(t, recipes)
}

// тест на поиск всех рецептов по юзеру
func TestRecipeRepo_FindByUser(t *testing.T) {
	ctx := context.Background()

	// создаем пользователя
	userID := primitive.NewObjectID()

	// создаем рецепты этого пользователя
	recipe1 := &domain.Recipe{
		Title:       "TestRecipe9",
		Description: "Test description",
		OwnerId:     userID,
		Ingredients: []domain.Ingredient{{Name: "Test", Quantity: "1"}},
		Steps:       []string{"Step 1"},
		Difficulty:  "easy",
	}
	err := testRepo.Create(ctx, recipe1)
	require.NoError(t, err)

	recipe2 := &domain.Recipe{
		Title:       "TestRecipe10",
		Description: "Test description",
		OwnerId:     userID,
		Ingredients: []domain.Ingredient{{Name: "Test", Quantity: "2"}},
		Steps:       []string{"Step 1"},
		Difficulty:  "medium",
	}
	err = testRepo.Create(ctx, recipe2)
	require.NoError(t, err)

	// ищем рецепты пользователя
	recipes, err := testRepo.FindByUser(ctx, userID)

	// сверяем че нашли
	assert.NoError(t, err)
	assert.NotNil(t, recipes)
	assert.GreaterOrEqual(t, len(recipes), 2)

	// проверяем что все рецепты принадлежат этому пользователю
	for _, recipe := range recipes {
		assert.Equal(t, userID, recipe.OwnerId)
	}
}

// тест на поиск всех рецептов юзера без рецептов
func TestRecipeRepo_FindByUser_NoRecipes(t *testing.T) {
	ctx := context.Background()

	// айди юзера без рецептов
	fakeUserID := primitive.NewObjectID()

	// получаем рецепты
	recipes, err := testRepo.FindByUser(ctx, fakeUserID)

	// не должны ничего найти
	assert.NoError(t, err)
	assert.Empty(t, recipes)
}

// тест на обновление рецепта
func TestRecipeRepo_Update(t *testing.T) {
	ctx := context.Background()

	// создаем тестовый рецепт
	recipe := createTestRecipe(t, "TestRecipe11")

	// данные для обновления
	newTitle := "NewTitle"
	newDescription := "NewDescript"
	newDifficulty := "hard"
	newPrepTime := 29
	newCookTime := 54

	// обновляем рецепт
	err := testRepo.Update(
		ctx,
		recipe.Id,
		recipe.OwnerId,
		&newTitle,
		&newDescription,
		nil,
		nil,
		&newPrepTime,
		&newCookTime,
		nil,
		&newDifficulty,
		nil,
	)
	assert.NoError(t, err)

	// получаем обновленный рецепт
	updated, err := testRepo.FindByID(ctx, recipe.Id)
	assert.NoError(t, err)

	// проверяем обновленные поля
	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, newDescription, updated.Description)
	assert.Equal(t, newDifficulty, updated.Difficulty)
	assert.Equal(t, newPrepTime, updated.PrepTime)
	assert.Equal(t, newCookTime, updated.CookTime)

	// проверяем необновленные поля
	assert.Equal(t, recipe.Ingredients, updated.Ingredients)
	assert.Equal(t, recipe.Steps, updated.Steps)
	assert.Equal(t, recipe.Servings, updated.Servings)
	assert.Equal(t, recipe.Tags, updated.Tags)
}

// тест на обновление рецепта без параметров
func TestRecipeRepo_Update_NothingToUpdate(t *testing.T) {
	ctx := context.Background()

	// создаем рецепт
	recipe := createTestRecipe(t, "TestRecipe12")

	// пытаемся обновить без параметров
	err := testRepo.Update(
		ctx,
		recipe.Id,
		recipe.OwnerId,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	// не должно получиться
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to update")
}

// тест на обновление рецепта без рецепта
func TestRecipeRepo_Update_NotFound(t *testing.T) {
	ctx := context.Background()

	// делаем фейк данные для рецепта
	fakeID := primitive.NewObjectID()
	fakeUserID := primitive.NewObjectID()
	newTitle := "TestRecipe13"

	// пытаемся обновить
	err := testRepo.Update(
		ctx,
		fakeID,
		fakeUserID,
		&newTitle,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	// не должно получиться
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// тест на обновление чужого рецепта
func TestRecipeRepo_Update_WrongOwner(t *testing.T) {
	ctx := context.Background()

	// создаем рецепт
	recipe := createTestRecipe(t, "Wrong Owner Update")

	// имитируем другого пользователя
	wrongUserID := primitive.NewObjectID()
	newTitle := "TestRecipe14"

	// пытаемся обновить рецепт
	err := testRepo.Update(
		ctx,
		recipe.Id,
		wrongUserID,
		&newTitle,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	// не должно получиться
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// тест на удаление рецепта
func TestRecipeRepo_Delete(t *testing.T) {
	ctx := context.Background()

	// создаем рецепт
	recipe := createTestRecipe(t, "TestRecipe15")

	// удаляем рецепт
	err := testRepo.Delete(ctx, recipe.Id, recipe.OwnerId)
	assert.NoError(t, err)

	// проверяем, что рецепт удален
	_, err = testRepo.FindByID(ctx, recipe.Id)
	assert.Error(t, err)
}

// тест на удаление рецепта без рецепта
func TestRecipeRepo_Delete_NotFound(t *testing.T) {
	ctx := context.Background()

	// фейк данные
	fakeID := primitive.NewObjectID()
	fakeUserID := primitive.NewObjectID()

	// пытаемся удалить
	err := testRepo.Delete(ctx, fakeID, fakeUserID)

	// не должно получиться
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// тест на удаление чужого рецепта
func TestRecipeRepo_Delete_WrongOwner(t *testing.T) {
	ctx := context.Background()

	// создаем рецепт
	recipe := createTestRecipe(t, "TestRecipe16")

	// фейк данные
	wrongUserID := primitive.NewObjectID()

	// пытаемся удалить рецепт
	err := testRepo.Delete(ctx, recipe.Id, wrongUserID)

	// нельзя
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
