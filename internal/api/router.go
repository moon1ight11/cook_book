package api

import "cook_book/internal/api/handlers"

type Router struct {
	recipeHandler handlers.RecipeHandlersInterface
}