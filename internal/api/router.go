package api

import (
	"cook_book/internal/api/handlers"
	"cook_book/internal/api/jwt"
	"cook_book/internal/api/middlewares"
	"cook_book/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// структура роутера
type Router struct {
	recipeHandler handlers.RecipeHandlersInterface
	authHandler   handlers.AuthHandlersInterface
	userHandler   handlers.UserHandlersInterface
	engine        *echo.Echo
}

// конструктор роутера
func NewRouter(
	recipehandler handlers.RecipeHandlersInterface,
	authHandler handlers.AuthHandlersInterface,
	userHandler handlers.UserHandlersInterface,
) *Router {
	return &Router{
		recipeHandler: recipehandler,
		authHandler:   authHandler,
		userHandler:   userHandler,
		engine:        echo.New(),
	}
}

// инициализация роутов
func (r *Router) Init(logger logger.Logger, jwtService jwt.JWTInterface) {
	// Middleware для корс
	r.engine.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins:     []string{"https://my-app.com", "http://localhost:3000"},
			AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowCredentials: true,
		},
	))

	// эндпоинт для проверки состояния
	r.engine.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// группировка роутов
	authGroup := r.engine.Group("/api/v1/auth")
	publicGroup := r.engine.Group("/api/v1/public")
	privateGroup := r.engine.Group("/api/v1/private")
	adminGroup := r.engine.Group("/api/v1/admin")

	// мидлвары на приватную группу для базовых юзеров
	privateGroup.Use(middlewares.Auth(jwtService, logger))
	privateGroup.Use(middlewares.RequireRole("Basic", "Admin"))

	// мидлвары на приватную группу для админов
	adminGroup.Use(middlewares.Auth(jwtService, logger))
	adminGroup.Use(middlewares.RequireRole("Admin"))

	// группа хэндлеров аутентификации
	authGroup.POST("/sign-up", r.authHandler.SignUp)
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	authGroup.POST("/token-refresh", r.authHandler.RefreshToken) 

	// группа публичных хэндлеров с рецептами
	publicGroup.GET("/recipes", r.recipeHandler.ListRecipes)
	publicGroup.GET("/recipes/:recipe_id", r.recipeHandler.GetRecipe)
	publicGroup.GET("/recipes-by-tags", r.recipeHandler.SearchByTags)
	publicGroup.GET("/recipes-by-user/:user_id", r.recipeHandler.SearchByUser)

	// приватная группа хэндлеров с рецептами
	privateGroup.POST("/recipes", r.recipeHandler.CreateRecipe)
	privateGroup.PATCH("/recipes/:recipe_id", r.recipeHandler.UpdateRecipe)
	privateGroup.DELETE("/recipes/:recipe_id", r.recipeHandler.DeleteRecipe)

	// приватная группа хэндлеров с пользователями
	privateGroup.GET("/users", r.userHandler.GetUser)
	privateGroup.PATCH("/users", r.userHandler.UpdateUser)
	privateGroup.DELETE("/users", r.userHandler.DeleteUser)
	privateGroup.POST("/sign-out", r.authHandler.SignOut)

	// админская группа
	adminGroup.POST("/users", r.userHandler.CreateAdmin)
	adminGroup.PATCH("/users/:user_id", r.userHandler.UpdateUserRole)
	adminGroup.DELETE("/users/:user_id", r.userHandler.DeleteAnyUser)
	adminGroup.DELETE("/recipe/:recipe_id", r.recipeHandler.DeleteAnyRecipe)
}

func (r *Router) GetEngine() *echo.Echo {
	return r.engine
}
