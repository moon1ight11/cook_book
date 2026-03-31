package main

import (
	"context"
	"cook_book/internal/api/handlers/recipehandlers"
	"cook_book/internal/config"
	"cook_book/internal/services/recipe"
	"cook_book/internal/storage"
	inmemcache "cook_book/internal/storage/cache/in_mem_cache"
	"cook_book/internal/storage/repo/recipes"
	"cook_book/pkg/logger"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// обработка паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	// инициализация конфига
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
	}

	// инициализация логгера
	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic(err)
	}

	// инициализация БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	// контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// инициализация кэша
	recipeCache := inmemcache.NewRecipeCache(cfg.Cache.TTL, cfg.Cache.CleanUpInterval)

	// инициализация зависимостей
	recipeRepo := recipes.NewRecipeMongoRepo(db, logger)
	recipeUseCase := recipe.NewRecipeUseCase(recipeRepo, recipeCache)
	recipeHandler := recipehandlers.NewRecipeHandler(recipeUseCase, logger)

	// Setup Echo
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	// API routes
	api := e.Group("/api/v1")

	// Recipe routes
	recipes := api.Group("/recipes")
	recipes.POST("", recipeHandler.CreateRecipe)
	recipes.GET("", recipeHandler.ListRecipes)
	recipes.GET("/search/tags", recipeHandler.SearchByTags)
	recipes.GET("/:id", recipeHandler.GetRecipe)
	recipes.PUT("/:id", recipeHandler.UpdateRecipe)
	recipes.DELETE("/:id", recipeHandler.DeleteRecipe)

	// Stats routes
	api.GET("/stats", recipeHandler.GetStats)

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown с таймаутом
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
	}
}
