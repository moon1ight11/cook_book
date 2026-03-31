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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
	}

	// валидируем поля рецепта
	if err := req.Validate(); err != nil {
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
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// возвращаем рецепт в виде апи-модели для ответа
	return c.JSON(http.StatusCreated, models.ToRecipeResponse(recipe))
}

// получение рецепта
func (h *RecipeHandler) GetRecipe(c echo.Context) error {
	// из параметров достаем id рецепта
	id := c.Param("id")

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// получаем рецепт
	recipe, err := h.useCase.GetRecipeByID(ctx, id)
	if err != nil {
		// если ошибка - не найдено
		if err == domain.ErrRecipeNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "recipe not found"})
		}
		// если какая-то другая ошибка
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	// возвращаем детальный рецепт
	return c.JSON(http.StatusOK, models.ToDetailedRecipeResponse(recipe))
}

// обновление рецепта
func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	// получаем id ид параметров
	id := c.Param("id")

	// приводим к типу для монго
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid recipe ID"})
	}

	// получаем модель для обновления
	var req models.UpdateRecipeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
	}

	// валидируем поля
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// создаем доменную модель для отправки в сервис
	recipe := &domain.Recipe{
		Id:          objectID,
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

	// обновление модели
	err = h.useCase.UpdateRecipe(ctx, recipe)
	if err != nil {
		// ошибка - рецепт не найден
		if err == domain.ErrRecipeNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "recipe not found"})
		}
		// другие ошибки
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// возвращаем сообщение, что модель обновлена
	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "recipe updated successfully"})
}

// удаление рецепта
func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	// получаем id из параметров
	id := c.Param("id")

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// обновление модели
	err := h.useCase.DeleteRecipe(ctx, id)
	if err != nil {
		// ошибка - рецепт не найден
		if err == domain.ErrRecipeNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "recipe not found"})
		}
		// другие ошибки
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

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
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// получаем массив рецептов
	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	// определяем количество страниц исходя из пагинации
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "at least one tag is required"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// получаем рецепты
	recipes, err := h.useCase.SearchByTags(ctx, tags)
	if err != nil {
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

// получение статистики
func (h *RecipeHandler) GetStats(c echo.Context) error {
	total, err := h.useCase.GetTotalCount(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, models.StatsResponse{TotalRecipes: total})
}
