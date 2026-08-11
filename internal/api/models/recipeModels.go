package models

import (
	"cook_book/internal/domain"
	"fmt"
	"strings"
	"time"
)

// создание рецепта для апи
type CreateRecipeRequest struct {
	Title       string              `json:"title" validate:"required,min=3,max=100"`
	Description string              `json:"description" validate:"max=500"`
	Ingredients []domain.Ingredient `json:"ingredients" validate:"required,min=1,dive"`
	Steps       []string            `json:"steps" validate:"required,min=1"`
	PrepTime    int                 `json:"prep_time" validate:"gte=0"`
	CookTime    int                 `json:"cook_time" validate:"gte=0"`
	Servings    int                 `json:"servings" validate:"gt=0"`
	Difficulty  string              `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Tags        []string            `json:"tags"`
}

// валидация созданного рецепта для апи
func (r *CreateRecipeRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return fmt.Errorf("Validate: empty title")
	}
	if len(r.Ingredients) == 0 {
		return fmt.Errorf("Validate: no ingredients")
	}
	if len(r.Steps) == 0 {
		return fmt.Errorf("Validate: no steps")
	}
	if r.Servings <= 0 {
		return fmt.Errorf("Validate: no servings")
	}
	return nil
}

// структура для обновления рецепта
type UpdateRecipe struct {
	Title       *string              `json:"title"`
	Description *string              `json:"description" validate:"max=500"`
	Ingredients *[]domain.Ingredient `json:"ingredients"`
	Steps       *[]string            `json:"steps"`
	PrepTime    *int                 `json:"prep_time"`
	CookTime    *int                 `json:"cook_time"`
	Servings    *int                 `json:"servings"`
	Difficulty  *string              `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Tags        *[]string            `json:"tags"`
}

// рецепт для ответа апи
type RecipeResponse struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PrepTime    int       `json:"prep_time"`
	CookTime    int       `json:"cook_time"`
	Servings    int       `json:"servings"`
	Difficulty  string    `json:"difficulty"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

// детальный рецепт для ответа апи
type DetailedRecipeResponse struct {
	ID          string              `json:"id"`
	OwnerID     string              `json:"owner_id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Ingredients []domain.Ingredient `json:"ingredients"`
	Steps       []string            `json:"steps"`
	PrepTime    int                 `json:"prep_time"`
	CookTime    int                 `json:"cook_time"`
	Servings    int                 `json:"servings"`
	Difficulty  string              `json:"difficulty"`
	Tags        []string            `json:"tags"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// ответ с пагинацией
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

// ответ с количеством рецептов
type StatsResponse struct {
	TotalRecipes int64 `json:"total_recipes"`
}

// перевод репозиторной модели в апи
func ToRecipeResponse(recipe *domain.Recipe) RecipeResponse {
	return RecipeResponse{
		ID:          recipe.Id.Hex(),
		OwnerID:     recipe.OwnerId.Hex(),
		Title:       recipe.Title,
		Description: recipe.Description,
		PrepTime:    recipe.PrepTime,
		CookTime:    recipe.CookTime,
		Servings:    recipe.Servings,
		Difficulty:  recipe.Difficulty,
		Tags:        recipe.Tags,
		CreatedAt:   recipe.CreatedAt,
	}
}

// перевод детальной репо-модели в апи
func ToDetailedRecipeResponse(recipe *domain.Recipe) DetailedRecipeResponse {
	return DetailedRecipeResponse{
		ID:          recipe.Id.Hex(),
		OwnerID:     recipe.OwnerId.Hex(),
		Title:       recipe.Title,
		Description: recipe.Description,
		Ingredients: recipe.Ingredients,
		Steps:       recipe.Steps,
		PrepTime:    recipe.PrepTime,
		CookTime:    recipe.CookTime,
		Servings:    recipe.Servings,
		Difficulty:  recipe.Difficulty,
		Tags:        recipe.Tags,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
	}
}
