package domain

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
		return fmt.Errorf("RecipeValidate: empty title")
	}
	if len(r.Ingredients) == 0 {
		return fmt.Errorf("RecipeValidate: no ingredients")
	}
	if len(r.Steps) == 0 {
		return fmt.Errorf("RecipeValidate: no steps")
	}
	if r.PrepTime < 0 {
		return fmt.Errorf("RecipeValidate: prep time must be positive")
	}
	if r.CookTime < 0 {
		return fmt.Errorf("RecipeValidate: cook time must be positive")
	}
	if r.Servings <= 0 {
		return fmt.Errorf("RecipeValidate: servings must be positive")
	}
	if r.Difficulty != "easy" && r.Difficulty != "medium" && r.Difficulty != "hard" {
		return fmt.Errorf("RecipeValidate: wrong difficulty")
	}
	return nil
}
