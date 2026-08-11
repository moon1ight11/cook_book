package api

import (
	"net/http"

	"cook_book/internal/api/handlers"
	"cook_book/internal/api/jwt"
	"cook_book/internal/api/middlewares"
	"cook_book/internal/metrics"
	"cook_book/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	recipeHandler handlers.RecipeHandlersInterface
	authHandler   handlers.AuthHandlersInterface
	userHandler   handlers.UserHandlersInterface
	engine        *echo.Echo
}

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

func (r *Router) Init(logger logger.Logger, jwtService jwt.JWTInterface, metrics *metrics.Metrics) {
	r.engine.Use(metrics.EchoMiddleware())

	r.engine.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowCredentials: true,
		},
	))

	r.engine.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	metrics.RegisterMetricsHandler(r.engine, "/metrics")

	authGroup := r.engine.Group("/api/v1/auth")
	publicGroup := r.engine.Group("/api/v1/public")
	privateGroup := r.engine.Group("/api/v1/private")
	adminGroup := r.engine.Group("/api/v1/admin")

	privateGroup.Use(middlewares.Auth(jwtService, logger))
	privateGroup.Use(middlewares.RequireRole("Basic", "Admin"))

	adminGroup.Use(middlewares.Auth(jwtService, logger))
	adminGroup.Use(middlewares.RequireRole("Admin"))

	authGroup.POST("/sign-up", r.authHandler.SignUp)
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	authGroup.POST("/token-refresh", r.authHandler.RefreshToken)

	publicGroup.GET("/recipes", r.recipeHandler.ListRecipes)
	publicGroup.GET("/recipes/:recipe_id", r.recipeHandler.GetRecipe)
	publicGroup.GET("/recipes-by-tags", r.recipeHandler.SearchByTags)
	publicGroup.GET("/recipes-by-user/:user_id", r.recipeHandler.SearchByUser)

	privateGroup.POST("/recipes", r.recipeHandler.CreateRecipe)
	privateGroup.PATCH("/recipes/:recipe_id", r.recipeHandler.UpdateRecipe)
	privateGroup.DELETE("/recipes/:recipe_id", r.recipeHandler.DeleteRecipe)

	privateGroup.GET("/users", r.userHandler.GetUser)
	privateGroup.PATCH("/users", r.userHandler.UpdateUser)
	privateGroup.DELETE("/users", r.userHandler.DeleteUser)
	privateGroup.POST("/sign-out", r.authHandler.SignOut)

	adminGroup.POST("/users", r.userHandler.CreateAdmin)
	adminGroup.PATCH("/users/:user_id", r.userHandler.UpdateUserRole)
	adminGroup.DELETE("/users/:user_id", r.userHandler.DeleteAnyUser)
	adminGroup.DELETE("/recipe/:recipe_id", r.recipeHandler.DeleteAnyRecipe)
}

func (r *Router) GetEngine() *echo.Echo {
	return r.engine
}
