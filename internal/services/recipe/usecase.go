package recipe

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"cook_book/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *recipeUseCase) CreateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	if err := recipe.RecipeValidate(); err != nil {
		return fmt.Errorf("CreateRecipe: %w", err)
	}

	if recipe.Difficulty == "" {
		recipe.Difficulty = "medium"
	}

	for i, tag := range recipe.Tags {
		recipe.Tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	if err := uc.repo.Create(ctx, recipe); err != nil {
		return fmt.Errorf("CreateRecipe: %w", err)
	}

	return nil
}

func (uc *recipeUseCase) UpdateRecipe(
	ctx context.Context,
	id primitive.ObjectID,
	userId primitive.ObjectID,
	title *string,
	description *string,
	ingredients *[]domain.Ingredient,
	steps *[]string,
	prepTime *int,
	cooktime *int,
	servings *int,
	difficulty *string,
	tags *[]string,
) error {
	if tags != nil {
		for i, tag := range *tags {
			(*tags)[i] = strings.ToLower(strings.TrimSpace(tag))
		}
	}

	if title != nil && strings.TrimSpace(*title) == "" {
		return fmt.Errorf("UpdateRecipe: title is empty")
	}

	if cooktime != nil && *cooktime <= 0 {
		return fmt.Errorf("UpdateRecipe: cook time must be positive")
	}

	if difficulty != nil && *difficulty != "hard" && *difficulty != "medium" && *difficulty != "easy" {
		return fmt.Errorf("UpdateRecipe: invalid difficulty: %s", *difficulty)
	}

	if err := uc.repo.Update(ctx, id, userId, title, description, ingredients, steps, prepTime, cooktime, servings, difficulty, tags); err != nil {
		return fmt.Errorf("UpdateRecipe: %w", err)
	}

	uc.cache.Delete(id.Hex())

	return nil
}

func (uc *recipeUseCase) DeleteRecipe(ctx context.Context, id string, userId string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("DeleteRecipe: %w", err)
	}

	objectUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("DeleteRecipe: %w", err)
	}

	if err = uc.repo.Delete(ctx, objectID, objectUserId); err != nil {
		return fmt.Errorf("DeleteRecipe: %w", err)
	}

	uc.cache.Delete(id)

	return nil
}

func (uc *recipeUseCase) GetRecipeByID(ctx context.Context, id string) (*domain.Recipe, error) {
	if cachedRecipe, found := uc.cache.Get(id); found {
		return cachedRecipe, nil
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("GetRecipeByID: %w", err)
	}

	recipe, err := uc.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("GetRecipeByID: %w", err)
	}

	uc.cache.Set(id, recipe)

	return recipe, nil
}

func (uc *recipeUseCase) GetRecipeByUser(ctx context.Context, userId string) ([]*domain.Recipe, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("GetRecipeByUser: %w", err)
	}

	recipes, err := uc.repo.FindByUser(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("GetRecipeByUser: %w", err)
	}

	return recipes, nil
}

func (uc *recipeUseCase) ListRecipes(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domain.Recipe, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	limit := int64(pageSize)
	offset := int64((page - 1) * pageSize)

	recipes, amount, err := uc.repo.FindAll(ctx, filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ListRecipes: %w", err)
	}

	return recipes, amount, nil
}

func (uc *recipeUseCase) SearchByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error) {
	if len(tags) == 0 {
		return []*domain.Recipe{}, fmt.Errorf("SearchByTags: no tags")
	}

	for i, tag := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	sortedTags := make([]string, len(tags))
	copy(sortedTags, tags)
	sort.Strings(sortedTags)
	cacheKey := strings.Join(sortedTags, ",")

	if cached, found := uc.cache.GetByTags(cacheKey); found {
		return cached, nil
	}

	recipes, err := uc.repo.FindByTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("SearchByTags: %w", err)
	}

	uc.cache.SetByTags(cacheKey, recipes)

	return recipes, nil
}
