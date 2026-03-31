package api

import (
	"cook_book/internal/api/handlers"
	"cook_book/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// структура роутера
type Router struct {
	recipeHandler handlers.RecipeHandlersInterface
	engine        *echo.Echo
}

// конструктор роутера
func NewRouter(recipehandler handlers.RecipeHandlersInterface) *Router {
	return &Router{
		recipeHandler: recipehandler,
		engine:        echo.New(),
	}
}

// инициализация роутов
func (r *Router) Init(logger logger.Logger) {
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
	privateGroup := r.engine.Group("/api/v1/private")
	publicGroup := r.engine.Group("/api/v1/public")
	// authGroup := r.engine.Group("/api/v1/auth")

	// создание рецепта
	privateGroup.POST("/recipes", r.recipeHandler.CreateRecipe)
	// изменение рецепта
	privateGroup.PATCH("/recipes", r.recipeHandler.UpdateRecipe)
	// удаление рецепта
	privateGroup.DELETE("/recipes", r.recipeHandler.DeleteRecipe)

	// просмотр всех рецептов с фильтрами и пагинацией
	publicGroup.GET("/recipes", r.recipeHandler.ListRecipes)
	// просмотр одного рецепта
	publicGroup.GET("/recipes/:recipe_id", r.recipeHandler.GetRecipe)
	// просмотр рецептов по тегам
	publicGroup.GET("/recipes", r.recipeHandler.SearchByTags)
}

func (r *Router) GetEngine() *echo.Echo {
	return r.engine
}
