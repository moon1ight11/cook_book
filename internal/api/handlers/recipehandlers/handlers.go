package recipehandlers

import (
	"context"
	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// создание рецепта
func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	var req models.CreateRecipeRequest
	// получаем рецепт с фронта
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Error in CreateRecipe:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
	}

	// валидируем поля рецепта
	if err := req.Validate(); err != nil {
		h.logger.Error("Error in CreateRecipe:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// создаем доменную модель для передачи в сервис
	recipe := &domain.Recipe{
		Title:       req.Title,
		Description: req.Description,
		Ingredients: req.Ingredients,
		Steps:       req.Steps,
		PrepTime:    req.PrepTime,
		CookTime:    req.CookTime,
		Servings:    req.Servings,
		Difficulty:  req.Difficulty,
		Tags:        req.Tags,
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// добавляем рецепт
	err := h.useCase.CreateRecipe(ctx, recipe)
	if err != nil {
		h.logger.Error("Error in CreateRecipe:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	h.logger.Info("Recipe created successful:", "_id", recipe.Id, "title", recipe.Title)

	// возвращаем рецепт в виде апи-модели для ответа
	return c.JSON(http.StatusCreated, models.ToRecipeResponse(recipe))
}

func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	// получаем id из параметров
	id := c.Param("recipe_id")

	// приводим в правильный формат
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.logger.Error("Error in UpdateRecipe:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "wrong id"})
	}

	// получаем id пользователя из контекста
	userId, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user id type")
	}

	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		h.logger.Error("Error in UpdateRecipe:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "wrong id"})
	}

	// получаем с фронта структуру с обновленным рецептом
	var req models.UpdateRecipe
	// получаем рецепт с фронта
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Error in CreateRecipe:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// обновляем рецепт
	err = h.useCase.UpdateRecipe(
		ctx,
		objectID,
		userID,
		req.Title,
		req.Description,
		req.Ingredients,
		req.Steps,
		req.PrepTime,
		req.CookTime,
		req.Servings,
		req.Difficulty,
		req.Tags,
	)
	if err != nil {
		h.logger.Error("Error in UpdateRecipe:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	h.logger.Info("Recipe updated successful:", "_id", objectID)

	// возвращаем сообщение о успешном удалении
	return c.JSON(http.StatusCreated, models.SuccessResponse{Message: "updated successful"})
}

// получение рецепта
func (h *RecipeHandler) GetRecipe(c echo.Context) error {
	// из параметров достаем id рецепта
	id := c.Param("recipe_id")

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// получаем рецепт
	recipe, err := h.useCase.GetRecipeByID(ctx, id)
	if err != nil {
		// если ошибка - не найдено
		if err == domain.ErrRecipeNotFound {
			h.logger.Error("Error in GetRecipe:", "error", err)
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "recipe not found"})
		}
		// если какая-то другая ошибка
		h.logger.Error("Error in GetRecipe:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	h.logger.Info("Recipe gets successful:", "_id", recipe.Id, "title", recipe.Title)

	// возвращаем детальный рецепт
	return c.JSON(http.StatusOK, models.ToDetailedRecipeResponse(recipe))
}

// удаление рецепта
func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	// получаем id из параметров
	id := c.Param("recipe_id")

	// получаем id пользователя из контекста
	userId, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user id type")
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// удаление рецепта
	err := h.useCase.DeleteRecipe(ctx, id, userId)
	if err != nil {
		// ошибка - рецепт не найден
		if err == domain.ErrRecipeNotFound {
			h.logger.Error("Error in DeleteRecipe:", "error", err)
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "recipe not found"})
		}
		// другие ошибки
		h.logger.Error("Error in DeleteRecipe:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	h.logger.Info("Recipe deleted successful:", "_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "recipe deleted successfully"})
}

// получение всех рецептов с фильрами и пагинацией
func (h *RecipeHandler) ListRecipes(c echo.Context) error {
	// получаем и валидируем данные для пагинации
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// получаем данные для фильтрации
	filters := make(map[string]interface{})
	if title := c.QueryParam("title"); title != "" {
		filters["title"] = title
	}
	if difficulty := c.QueryParam("difficulty"); difficulty != "" {
		filters["difficulty"] = difficulty
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// получаем рецепты и коничество
	recipes, total, err := h.useCase.ListRecipes(ctx, filters, page, pageSize)
	if err != nil {
		h.logger.Error("Error in ListRecipes:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// получаем массив рецептов
	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	// определяем количество страниц исходя из пагинации
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	h.logger.Info("Recipes list get successful:", "total", total)

	// возвращаем ответ с пагинацией
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// получение рецептов по тегам
func (h *RecipeHandler) SearchByTags(c echo.Context) error {
	// ищем теги в квери параметрах
	tags := c.QueryParams()["tags"]
	if len(tags) == 0 {
		h.logger.Error("Error in SearchByTags:", "error", "len(tags) = 0")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "at least one tag is required"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// получаем рецепты
	recipes, err := h.useCase.SearchByTags(ctx, tags)
	if err != nil {
		h.logger.Error("Error in SearchByTags:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// создаем модели для ответа
	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	h.logger.Info("Recipes by tags get successful:", "total", len(recipes)-1)

	// возвращаем массив рецептов
	return c.JSON(http.StatusOK, responses)
}

// поиск рецептов по пользователю
func (h *RecipeHandler) SearchByUser(c echo.Context) error {
	// из параметров достаем id пользователя
	id := c.Param("recipe_id")

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	recipes, err := h.useCase.GetRecipeByUser(ctx, id)
	if err != nil {
		h.logger.Error("Error in SearchByTags:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// создаем модели для ответа
	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	// возвращаем массив рецептов
	return c.JSON(http.StatusOK, responses)
}
