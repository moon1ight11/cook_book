package recipe

import (
	"context"
	"cook_book/internal/domain"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// создание рецепта
func (uc *recipeUseCase) CreateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	// валидируем входящий рецепт
	if err := recipe.RecipeValidate(); err != nil {
		return fmt.Errorf("error in CreateRecipe: %w", err)
	}

	// если сложность не установлена - ставим среднюю по дефолту
	if recipe.Difficulty == "" {
		recipe.Difficulty = "medium"
	}

	// делаем, чтобы теги всегда были в нижнем регистре
	for i, tag := range recipe.Tags {
		recipe.Tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	// создаем запись
	err := uc.repo.Create(ctx, recipe)
	if err != nil {
		return fmt.Errorf("error in CreateRecipe: %w", err)
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
	// делаем, чтобы теги всегда были в нижнем регистре
	if tags != nil {
		for i, tag := range *tags {
			(*tags)[i] = strings.ToLower(strings.TrimSpace(tag))
		}
	}

	// валидируем обязательные поля
	if title != nil {
		if strings.TrimSpace(*title) == "" {
			return fmt.Errorf("error in UpdateRecipe: title is empty")
		}
	}

	if cooktime != nil {
		if *cooktime <= 0 {
			return fmt.Errorf("error in UpdateRecipe: cook time must be positive")
		}
	}

	if difficulty != nil {
		if *difficulty != "hard" && *difficulty != "medium" && *difficulty != "easy" {
			return fmt.Errorf("error in UpdateRecipe: difficulty field is not valid: %s", *difficulty)
		}
	}

	// обновляем
	err := uc.repo.Update(ctx, id, userId, title, description, ingredients, steps, prepTime, cooktime, servings, difficulty, tags)
	if err != nil {
		return fmt.Errorf("error in UpdateRecipe: %w", err)
	}

	// удаляем запись из кэша
	uc.cache.Delete(id.Hex())

	return nil
}

// удаление записи
func (uc *recipeUseCase) DeleteRecipe(ctx context.Context, id string, userId string) error {
	// получаем id рецепта
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.NewRecipeError("invalid recipe ID")
	}

	objectUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return domain.NewRecipeError("invalid recipe ID")
	}

	// удаляем из бд
	err = uc.repo.Delete(ctx, objectID, objectUserId)
	if err != nil {
		return fmt.Errorf("error in DeleteRecipe: %w", err)
	}

	// удаляем запись из кэша
	uc.cache.Delete(id)

	return nil
}

// получение рецепта по id
func (uc *recipeUseCase) GetRecipeByID(ctx context.Context, id string) (*domain.Recipe, error) {
	// пробуем найти в кэше
	if cachedRecipe, found := uc.cache.Get(id); found {
		return cachedRecipe, nil
	}

	// в кэше не нашли - переводим формат id
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.NewRecipeError("invalid recipe ID")
	}

	// ищем в бд
	recipe, err := uc.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("error in GetRecipeById: %w", err)
	}

	// записываем в кэш
	uc.cache.Set(id, recipe)

	return recipe, nil
}

// получение списка рецептов конкретного пользователя
func (uc *recipeUseCase) GetRecipeByUser(ctx context.Context, userId string) ([]*domain.Recipe, error) {
	//переводим формат id
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, domain.NewRecipeError("invalid user ID")
	}

	// ищем в бд
	recipes, err := uc.repo.FindByUser(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("error in GetRecipeByUser: %w", err)
	}

	return recipes, nil
}

// получаем список рецептов с фильтрами
func (uc *recipeUseCase) ListRecipes(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domain.Recipe, int64, error) {
	// валидируем фильтры
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// устанавливаем размер и номер страницы
	limit := int64(pageSize)
	offset := int64((page - 1) * pageSize)

	// получаем из бд
	recipies, amount, err := uc.repo.FindAll(ctx, filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error in ListRecipes: %w", err)
	}

	return recipies, amount, nil
}

// поиск рецептов по тегу
func (uc *recipeUseCase) SearchByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error) {
	// проверяем, чтобы теги были
	if len(tags) == 0 {
		return []*domain.Recipe{}, nil
	}

	// преводим теги в нижний регистр
	for i, tag := range tags {
		tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	// соединяем теги в строку
	cacheKey := strings.Join(tags, ",")

	// ищем в кэше
	if cached, found := uc.cache.GetByTags(cacheKey); found {
		return cached, nil
	}

	// в кэше не нашли - идем в бд
	recipes, err := uc.repo.FindByTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("error in SearchByTags: %w", err)
	}

	// записываем в кэш
	uc.cache.SetByTags(cacheKey, recipes)

	return recipes, nil
}
