package recipehandlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"cook_book/internal/metrics"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *RecipeHandler) CreateRecipe(c echo.Context) error {
	h.metrics.RecordOperation("CreateRecipe")

	userId, ok := c.Get("user_id").(string)
	if !ok {
		h.metrics.RecordError(string(metrics.ErrForbidden), "CreateRecipe")
		h.logger.Error("CreateRecipe: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	ownerId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrForbidden), "CreateRecipe")
		h.logger.Error("CreateRecipe: invalid user_id", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	var req models.CreateRecipeRequest
	if err := c.Bind(&req); err != nil {
		h.metrics.RecordError(string(metrics.ErrBadRequest), "CreateRecipe")
		h.logger.Error("CreateRecipe: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	if err := req.Validate(); err != nil {
		h.metrics.RecordError(string(metrics.ErrForbidden), "CreateRecipe")
		h.logger.Error("CreateRecipe: validation failed", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	recipe := &domain.Recipe{
		Title:       req.Title,
		OwnerId:     ownerId,
		Description: req.Description,
		Ingredients: req.Ingredients,
		Steps:       req.Steps,
		PrepTime:    req.PrepTime,
		CookTime:    req.CookTime,
		Servings:    req.Servings,
		Difficulty:  req.Difficulty,
		Tags:        req.Tags,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err = h.useCase.CreateRecipe(ctx, recipe); err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "CreateRecipe")
		h.logger.Error("CreateRecipe: create recipe", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	h.logger.Info("Recipe created", "_id", recipe.Id, "title", recipe.Title)

	return c.JSON(http.StatusCreated, models.ToRecipeResponse(recipe))
}

func (h *RecipeHandler) UpdateRecipe(c echo.Context) error {
	h.metrics.RecordOperation("UpdateRecipe")

	id := c.Param("recipe_id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateRecipe")
		h.logger.Error("UpdateRecipe: invalid recipe_id", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	userId, ok := c.Get("user_id").(string)
	if !ok {
		h.metrics.RecordError(string(metrics.ErrForbidden), "UpdateRecipe")
		h.logger.Error("UpdateRecipe: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrForbidden), "UpdateRecipe")
		h.logger.Error("UpdateRecipe: invalid user_id", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	var req models.UpdateRecipe
	if err := c.Bind(&req); err != nil {
		h.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateRecipe")
		h.logger.Error("UpdateRecipe: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

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
		h.metrics.RecordError(string(metrics.ErrInternal), "UpdateRecipe")
		h.logger.Error("UpdateRecipe: update recipe", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	h.logger.Info("Recipe updated", "_id", objectID)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "updated successful"})
}

func (h *RecipeHandler) GetRecipe(c echo.Context) error {
	h.metrics.RecordOperation("GetRecipe")

	id := c.Param("recipe_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	recipe, err := h.useCase.GetRecipeByID(ctx, id)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "GetRecipe")
		h.logger.Error("GetRecipe: get recipe", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	return c.JSON(http.StatusOK, models.ToDetailedRecipeResponse(recipe))
}

func (h *RecipeHandler) DeleteRecipe(c echo.Context) error {
	h.metrics.RecordOperation("DeleteRecipe")

	id := c.Param("recipe_id")

	userId, ok := c.Get("user_id").(string)
	if !ok {
		h.metrics.RecordError(string(metrics.ErrForbidden), "DeleteRecipe")
		h.logger.Error("DeleteRecipe: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err := h.useCase.DeleteRecipe(ctx, id, userId); err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "DeleteRecipe")
		h.logger.Error("DeleteRecipe: delete recipe", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	h.logger.Info("Recipe deleted", "_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "recipe deleted successfully"})
}

func (h *RecipeHandler) ListRecipes(c echo.Context) error {
	h.metrics.RecordOperation("ListRecipes")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	filters := make(map[string]interface{})
	if title := c.QueryParam("title"); title != "" {
		filters["title"] = title
	}
	if difficulty := c.QueryParam("difficulty"); difficulty != "" {
		filters["difficulty"] = difficulty
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	recipes, total, err := h.useCase.ListRecipes(ctx, filters, page, pageSize)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "ListRecipes")
		h.logger.Error("ListRecipes: list recipes", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	h.logger.Info("Recipes listed", "total", total)

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *RecipeHandler) SearchByTags(c echo.Context) error {
	h.metrics.RecordOperation("SearchByTags")

	tags := c.QueryParams()["tags"]
	if len(tags) == 0 {
		h.metrics.RecordError(string(metrics.ErrBadRequest), "SearchByTags")
		h.logger.Error("SearchByTags: no tags provided")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	recipes, err := h.useCase.SearchByTags(ctx, tags)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "SearchByTags")
		h.logger.Error("SearchByTags: search by tags", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *RecipeHandler) SearchByUser(c echo.Context) error {
	h.metrics.RecordOperation("SearchByUser")

	id := c.Param("user_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	recipes, err := h.useCase.GetRecipeByUser(ctx, id)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "SearchByUser")
		h.logger.Error("SearchByUser: get recipe by user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	var responses []models.RecipeResponse
	for _, recipe := range recipes {
		responses = append(responses, models.ToRecipeResponse(recipe))
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *RecipeHandler) DeleteAnyRecipe(c echo.Context) error {
	h.metrics.RecordOperation("DeleteRecipe_ADMIN")

	id := c.Param("recipe_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	recipe, err := h.useCase.GetRecipeByID(ctx, id)
	if err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "DeleteRecipe_ADMIN")
		h.logger.Error("DeleteAnyRecipe: get recipe", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	if err = h.useCase.DeleteRecipe(ctx, id, recipe.OwnerId.Hex()); err != nil {
		h.metrics.RecordError(string(metrics.ErrInternal), "DeleteRecipe_ADMIN")
		h.logger.Error("DeleteAnyRecipe: delete recipe", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	h.logger.Info("Recipe deleted by admin", "_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "recipe deleted successfully"})
}
