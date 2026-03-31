package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// структура самого рецепта
type Recipe struct {
	Id primitive.ObjectID
	Title string
	Description string
	Ingredients []Ingredient
	Steps []string
	PrepTime int
	CookTime int
	Servings int
	Difficulty string
	Tags []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// структура ингредиентов
type Ingredient struct {
	Name string
	Quantity string
	Unit string
}

// валидация рецепта
func (r *Recipe) RecipeValidate() error {
    if r.Title == "" {
        return ErrEmptyTitle
    }
    if len(r.Ingredients) == 0 {
        return ErrNoIngredients
    }
    if len(r.Steps) == 0 {
        return ErrNoSteps
    }
    if r.PrepTime < 0 {
        return ErrInvalidPrepTime
    }
    if r.CookTime < 0 {
        return ErrInvalidCookTime
    }
    if r.Servings <= 0 {
        return ErrInvalidServings
    }
    return nil
}

// ошибки рецептов
var (
    ErrEmptyTitle       = NewRecipeError("recipe title is required")
    ErrNoIngredients    = NewRecipeError("recipe must have at least one ingredient")
    ErrNoSteps          = NewRecipeError("recipe must have at least one step")
    ErrRecipeNotFound   = NewRecipeError("recipe not found")
    ErrInvalidPrepTime  = NewRecipeError("prep time cannot be negative")
    ErrInvalidCookTime  = NewRecipeError("cook time cannot be negative")
    ErrInvalidServings  = NewRecipeError("servings must be greater than 0")
)

// структура ошибки рецептов
type RecipeError struct {
    Message string
}

func (e *RecipeError) Error() string {
    return e.Message
}

// конструктор ошибки рецептов
func NewRecipeError(msg string) *RecipeError {
    return &RecipeError{Message: msg}
}