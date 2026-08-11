package recipe

import (
	"context"
	"testing"

	"cook_book/internal/domain"
	"cook_book/internal/services/recipe/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTest() (*recipeUseCase, *mocks.MockRecipeRepository, *mocks.MockRecipeCache) {
	mockRepo := new(mocks.MockRecipeRepository)
	mockCache := new(mocks.MockRecipeCache)

	usecase := &recipeUseCase{
		repo:  mockRepo,
		cache: mockCache,
	}

	return usecase, mockRepo, mockCache
}

func createTestRecipe() *domain.Recipe {
	return &domain.Recipe{
		Id:          primitive.NewObjectID(),
		Title:       "TestRecipe",
		Description: "Test Description",
		OwnerId:     primitive.NewObjectID(),
		Ingredients: []domain.Ingredient{
			{Name: "Ingredient 1", Quantity: "100g"},
		},
		Steps:      []string{"Step 1", "Step 2"},
		PrepTime:   10,
		CookTime:   20,
		Servings:   4,
		Difficulty: "easy",
		Tags:       []string{"Test", "Simple"},
	}
}

func TestCreateRecipe_Success(t *testing.T) {
	uc, mockRepo, _ := setupTest()
	ctx := context.Background()

	recipe := createTestRecipe()
	recipe.Tags = []string{"  TEST  ", "Simple"}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Recipe")).Return(nil)

	err := uc.CreateRecipe(ctx, recipe)

	assert.NoError(t, err)
	assert.Equal(t, "easy", recipe.Difficulty)
	assert.Equal(t, []string{"test", "simple"}, recipe.Tags)
	mockRepo.AssertExpectations(t)
}

func TestCreateRecipe_InvalidRecipe(t *testing.T) {
	uc, mockRepo, _ := setupTest()
	ctx := context.Background()

	recipe := &domain.Recipe{
		Title: "",
	}

	err := uc.CreateRecipe(ctx, recipe)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreateRecipe")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUpdateRecipe_Success(t *testing.T) {
	uc, mockRepo, mockCache := setupTest()
	ctx := context.Background()

	recipe := createTestRecipe()

	var nilStr *string
	var nilIngredients *[]domain.Ingredient
	var nilSteps *[]string
	var nilInt *int
	var nilTags *[]string

	newDifficulty := "hard"

	mockRepo.On("Update", ctx, recipe.Id, recipe.OwnerId,
		nilStr, nilStr, nilIngredients, nilSteps,
		nilInt, nilInt, nilInt,
		&newDifficulty, nilTags).Return(nil)
	mockCache.On("Delete", recipe.Id.Hex()).Return()

	err := uc.UpdateRecipe(ctx, recipe.Id, recipe.OwnerId,
		nil, nil, nil, nil, nil, nil, nil, &newDifficulty, nil)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUpdateRecipe_EmptyTitle(t *testing.T) {
	uc, _, _ := setupTest()
	ctx := context.Background()

	emptyTitle := ""

	err := uc.UpdateRecipe(ctx, primitive.NewObjectID(), primitive.NewObjectID(),
		&emptyTitle, nil, nil, nil, nil, nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is empty")
}

func TestDeleteRecipe_Success(t *testing.T) {
	uc, mockRepo, mockCache := setupTest()
	ctx := context.Background()

	recipe := createTestRecipe()

	mockRepo.On("Delete", ctx, recipe.Id, recipe.OwnerId).Return(nil)
	mockCache.On("Delete", recipe.Id.Hex()).Return()

	err := uc.DeleteRecipe(ctx, recipe.Id.Hex(), recipe.OwnerId.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestDeleteRecipe_InvalidID(t *testing.T) {
	uc, _, _ := setupTest()
	ctx := context.Background()

	err := uc.DeleteRecipe(ctx, "IDNOTFOUND", primitive.NewObjectID().Hex())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DeleteRecipe")
}

func TestGetRecipeByID_CacheMiss_RepoHit(t *testing.T) {
	uc, mockRepo, mockCache := setupTest()
	ctx := context.Background()

	recipe := createTestRecipe()

	mockCache.On("Get", recipe.Id.Hex()).Return(nil, false)
	mockRepo.On("FindByID", ctx, recipe.Id).Return(recipe, nil)
	mockCache.On("Set", recipe.Id.Hex(), recipe).Return()

	result, err := uc.GetRecipeByID(ctx, recipe.Id.Hex())

	assert.NoError(t, err)
	assert.Equal(t, recipe.Id, result.Id)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetRecipeByID_InvalidID(t *testing.T) {
	uc, _, mockCache := setupTest()
	ctx := context.Background()

	mockCache.On("Get", "IDNOTFOUND").Return(nil, false)

	_, err := uc.GetRecipeByID(ctx, "IDNOTFOUND")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GetRecipeByID")
}

func TestGetRecipeByUser_Success(t *testing.T) {
	uc, mockRepo, _ := setupTest()
	ctx := context.Background()

	userID := primitive.NewObjectID()
	recipes := []*domain.Recipe{
		createTestRecipe(),
		createTestRecipe(),
	}

	mockRepo.On("FindByUser", ctx, userID).Return(recipes, nil)

	result, err := uc.GetRecipeByUser(ctx, userID.Hex())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetRecipeByUser_InvalidID(t *testing.T) {
	uc, _, _ := setupTest()
	ctx := context.Background()

	_, err := uc.GetRecipeByUser(ctx, "IDNOTFOUND")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GetRecipeByUser")
}

func TestListRecipes_Success(t *testing.T) {
	uc, mockRepo, _ := setupTest()
	ctx := context.Background()

	recipes := []*domain.Recipe{
		createTestRecipe(),
		createTestRecipe(),
	}
	total := int64(10)

	filters := map[string]interface{}{"difficulty": "easy"}

	mockRepo.On("FindAll", ctx, filters, int64(10), int64(0)).Return(recipes, total, nil)

	result, count, err := uc.ListRecipes(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, total, count)
	mockRepo.AssertExpectations(t)
}

func TestSearchByTags_Success(t *testing.T) {
	uc, mockRepo, mockCache := setupTest()
	ctx := context.Background()

	tags := []string{"TAG1", "  Tag2  "}
	recipes := []*domain.Recipe{
		createTestRecipe(),
	}

	cacheKey := "tag1,tag2"

	mockCache.On("GetByTags", cacheKey).Return(nil, false)
	mockRepo.On("FindByTags", ctx, []string{"tag1", "tag2"}).Return(recipes, nil)
	mockCache.On("SetByTags", cacheKey, recipes).Return()

	result, err := uc.SearchByTags(ctx, tags)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestSearchByTags_NoTags(t *testing.T) {
	uc, _, _ := setupTest()
	ctx := context.Background()

	_, err := uc.SearchByTags(ctx, []string{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no tags")
}
