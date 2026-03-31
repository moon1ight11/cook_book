package recipe

import (
	"context"
	"cook_book/internal/domain"
	"strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// создание рецепта
func (uc *recipeUseCase) CreateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	// валидируем входящий рецепт
	if err := recipe.RecipeValidate(); err != nil {
		return err
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
		return err
	}

	return nil
}

// обновление рецепта
func (uc *recipeUseCase) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	// валидируем входные данные
	if err := recipe.RecipeValidate(); err != nil {
		return err
	}

	// удаляем запись из кэша
	uc.cache.Delete(recipe.Id.Hex())

	// обновляем запись в бд
	err := uc.repo.Update(ctx, recipe)
	if err != nil {
		return err
	}

	return nil
}

// удаление записи
func (uc *recipeUseCase) DeleteRecipe(ctx context.Context, id string) error {
	// получаем id рецепта
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.NewRecipeError("invalid recipe ID")
	}

	// удаляем запись из кэша
	uc.cache.Delete(id)

	// удаляем из бд
	err = uc.repo.Delete(ctx, objectID)
	if err != nil {
		return err
	}

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
		return nil, err
	}

	// записываем в кэш
	uc.cache.Set(id, recipe)

	return recipe, nil
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
		return nil, 0, err
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
		return nil, err
	}

	// записываем в кэш
	uc.cache.SetByTags(cacheKey, recipes)

	return recipes, nil
}

// получение количества рецептов в сервисе
func (uc *recipeUseCase) GetTotalCount(ctx context.Context) (int64, error) {
	// идем в бд за количеством
	_, total, err := uc.repo.FindAll(ctx, map[string]interface{}{}, 1, 0)
	if err != nil {
		return 0, err
	}
	return total, nil
}
