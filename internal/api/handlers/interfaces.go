package handlers

import "github.com/labstack/echo/v4"

type RecipeHandlersInterface interface {
	CreateRecipe(c echo.Context) error
	GetRecipe(c echo.Context) error
	UpdateRecipe(c echo.Context) error
	DeleteRecipe(c echo.Context) error
	ListRecipes(c echo.Context) error
	SearchByTags(c echo.Context) error
}