package domain

import (
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// структура самого рецепта
type Recipe struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId     primitive.ObjectID `bson:"owner_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Ingredients []Ingredient       `bson:"ingredients"`
	Steps       []string           `bson:"steps"`
	PrepTime    int                `bson:"prep_time"`
	CookTime    int                `bson:"cook_time"`
	Servings    int                `bson:"servings"`
	Difficulty  string             `bson:"difficulty"`
	Tags        []string           `bson:"tags"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// структура ингредиентов
type Ingredient struct {
	Name     string `bson:"name"`
	Quantity string `bson:"quantity"`
	Unit     string `bson:"unit"`
}

// валидация рецепта
func (r *Recipe) RecipeValidate() error {
	if r.Title == "" {
		return fmt.Errorf("error in RecipeValidate: empty title")
	}
	if len(r.Ingredients) == 0 {
		return fmt.Errorf("error in RecipeValidate: no ingredients")
	}
	if len(r.Steps) == 0 {
		return fmt.Errorf("error in RecipeValidate: no steps")
	}
	if r.PrepTime < 0 {
		return fmt.Errorf("error in RecipeValidate: prep time must be positive")
	}
	if r.CookTime < 0 {
		return fmt.Errorf("error in RecipeValidate: cook time must be positive")
	}
	if r.Servings <= 0 {
		return fmt.Errorf("error in RecipeValidate: num of servings must be positive")
	}
	return nil
}