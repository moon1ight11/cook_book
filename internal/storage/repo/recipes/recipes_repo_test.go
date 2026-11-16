package recipes

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"cook_book/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testRepo *recipeRepo
var testClient *mongo.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongoURI := "mongodb://admin:123@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	testClient = client
	testDB := client.Database("cook_book_test")
	testCollection := testDB.Collection("recipes_test")

	testRepo = &recipeRepo{
		collection: testCollection,
		logger:     slog.Default(),
	}

	code := m.Run()

	testCollection.Drop(ctx)
	client.Disconnect(ctx)

	os.Exit(code)
}

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

func TestRecipeRepo_Create(t *testing.T) {
	ctx := context.Background()

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

	err := testRepo.Create(ctx, recipe)

	assert.NoError(t, err)
	assert.NotEmpty(t, recipe.Id)
	assert.False(t, recipe.CreatedAt.IsZero())
}

func TestRecipeRepo_FindByID(t *testing.T) {
	ctx := context.Background()

	recipe := createTestRecipe(t, "TestRecipe1")

	found, err := testRepo.FindByID(ctx, recipe.Id)

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, recipe.Title, found.Title)
	assert.Equal(t, recipe.Description, found.Description)
	assert.Equal(t, recipe.OwnerId, found.OwnerId)
}

func TestRecipeRepo_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()

	fakeID := primitive.NewObjectID()

	_, err := testRepo.FindByID(ctx, fakeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRecipeRepo_FindAll(t *testing.T) {
	ctx := context.Background()

	createTestRecipe(t, "TestRecipe3")
	createTestRecipe(t, "TestRecipe4")
	createTestRecipe(t, "TestRecipe5")

	recipes, total, err := testRepo.FindAll(ctx, map[string]interface{}{}, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, recipes)
	assert.GreaterOrEqual(t, total, int64(3))

	filter := map[string]interface{}{
		"title": "TestRecipe3",
	}
	recipes, total, err = testRepo.FindAll(ctx, filter, 10, 0)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	if len(recipes) > 0 {
		assert.Contains(t, recipes[0].Title, "TestRecipe3")
	}

	recipes, total, err = testRepo.FindAll(ctx, map[string]interface{}{}, 2, 0)
	assert.NoError(t, err)
	assert.LessOrEqual(t, int64(len(recipes)), int64(2))
	assert.GreaterOrEqual(t, total, int64(3))
}

func TestRecipeRepo_FindAll_WithDifficultyFilter(t *testing.T) {
	ctx := context.Background()

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

func TestRecipeRepo_FindByTags(t *testing.T) {
	ctx := context.Background()

	createTestRecipe(t, "TestRecipe7")
	createTestRecipe(t, "TestRecipe8")

	recipes, err := testRepo.FindByTags(ctx, []string{"test"})

	assert.NoError(t, err)
	assert.NotNil(t, recipes)
	assert.GreaterOrEqual(t, len(recipes), 2)
}

func TestRecipeRepo_FindByTags_NotFound(t *testing.T) {
	ctx := context.Background()

	recipes, err := testRepo.FindByTags(ctx, []string{"TAGNOTFOUND"})

	assert.NoError(t, err)
	assert.Empty(t, recipes)
}

func TestRecipeRepo_FindByUser(t *testing.T) {
	ctx := context.Background()

	userID := primitive.NewObjectID()

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

	recipes, err := testRepo.FindByUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, recipes)
	assert.GreaterOrEqual(t, len(recipes), 2)

	for _, recipe := range recipes {
		assert.Equal(t, userID, recipe.OwnerId)
	}
}

func TestRecipeRepo_FindByUser_NoRecipes(t *testing.T) {
	ctx := context.Background()

	fakeUserID := primitive.NewObjectID()

	recipes, err := testRepo.FindByUser(ctx, fakeUserID)

	assert.NoError(t, err)
	assert.Empty(t, recipes)
}

func TestRecipeRepo_Update(t *testing.T) {
	ctx := context.Background()

	recipe := createTestRecipe(t, "TestRecipe11")

	newTitle := "NewTitle"
	newDescription := "NewDescript"
	newDifficulty := "hard"
	newPrepTime := 29
	newCookTime := 54

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

	updated, err := testRepo.FindByID(ctx, recipe.Id)
	assert.NoError(t, err)

	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, newDescription, updated.Description)
	assert.Equal(t, newDifficulty, updated.Difficulty)
	assert.Equal(t, newPrepTime, updated.PrepTime)
	assert.Equal(t, newCookTime, updated.CookTime)

	assert.Equal(t, recipe.Ingredients, updated.Ingredients)
	assert.Equal(t, recipe.Steps, updated.Steps)
	assert.Equal(t, recipe.Servings, updated.Servings)
	assert.Equal(t, recipe.Tags, updated.Tags)
}

func TestRecipeRepo_Update_NothingToUpdate(t *testing.T) {
	ctx := context.Background()

	recipe := createTestRecipe(t, "TestRecipe12")

	err := testRepo.Update(
		ctx,
		recipe.Id,
		recipe.OwnerId,
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to update")
}

func TestRecipeRepo_Update_NotFound(t *testing.T) {
	ctx := context.Background()

	fakeID := primitive.NewObjectID()
	fakeUserID := primitive.NewObjectID()
	newTitle := "TestRecipe13"

	err := testRepo.Update(
		ctx,
		fakeID,
		fakeUserID,
		&newTitle,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRecipeRepo_Update_WrongOwner(t *testing.T) {
	ctx := context.Background()

	recipe := createTestRecipe(t, "Wrong Owner Update")

	wrongUserID := primitive.NewObjectID()
	newTitle := "TestRecipe14"

	err := testRepo.Update(
		ctx,
		recipe.Id,
		wrongUserID,
		&newTitle,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRecipeRepo_Delete(t *testing.T) {
	ctx := context.Background()

	recipe := createTestRecipe(t, "TestRecipe15")

	err := testRepo.Delete(ctx, recipe.Id, recipe.OwnerId)
	assert.NoError(t, err)

	_, err = testRepo.FindByID(ctx, recipe.Id)
	assert.Error(t, err)
}

func TestRecipeRepo_Delete_NotFound(t *testing.T) {
	ctx := context.Background()

	fakeID := primitive.NewObjectID()
	fakeUserID := primitive.NewObjectID()

	err := testRepo.Delete(ctx, fakeID, fakeUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRecipeRepo_Delete_WrongOwner(t *testing.T) {
	ctx := context.Background()

	recipe := createTestRecipe(t, "TestRecipe16")

	wrongUserID := primitive.NewObjectID()

	err := testRepo.Delete(ctx, recipe.Id, wrongUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
