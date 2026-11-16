package handlers

import "github.com/labstack/echo/v4"

type RecipeHandlersInterface interface {
	CreateRecipe(c echo.Context) error
	UpdateRecipe(c echo.Context) error
	GetRecipe(c echo.Context) error
	SearchByUser(c echo.Context) error
	DeleteRecipe(c echo.Context) error
	DeleteAnyRecipe(c echo.Context) error
	ListRecipes(c echo.Context) error
	SearchByTags(c echo.Context) error
}

type AuthHandlersInterface interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	RefreshToken(c echo.Context) error
	SignOut(c echo.Context) error
}

type UserHandlersInterface interface {
	CreateAdmin(c echo.Context) error
	GetUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	UpdateUserRole(c echo.Context) error
	DeleteUser(c echo.Context) error
	DeleteAnyUser(c echo.Context) error
}
