package main

import (
	"context"
	"cook_book/internal/api"
	"cook_book/internal/api/handlers/recipehandlers"
	"cook_book/internal/config"
	"cook_book/internal/services/recipe"
	"cook_book/internal/storage"
	inmemcache "cook_book/internal/storage/cache/in_mem_cache"
	"cook_book/internal/storage/repo/recipes"
	"cook_book/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Error("Failed to connect DB", "error:", err)
		os.Exit(1)
	}

	// инициализация кэша
	recipeCache := inmemcache.NewRecipeCache(cfg.Cache.TTL, cfg.Cache.CleanUpInterval)

	// инициализация зависимостей
	recipeRepo := recipes.NewRecipeMongoRepo(db, logger)
	recipeUseCase := recipe.NewRecipeUseCase(recipeRepo, recipeCache)
	recipeHandler := recipehandlers.NewRecipeHandler(recipeUseCase, logger)

	// инициализация роутера
	router := api.NewRouter(recipeHandler)

	// инициализация роутов
	router.Init(logger)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router.GetEngine(),
	}

	// создаем каналы для сигналов завершения и ошибки сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	// запуск сервера
	go func() {
		logger.Info("Server is starting on port :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to run server:", "error", err)
			serverError <- err
		}
	}()

	// ждем сигналы
	select {
	case <-quit:
		// если поступил сигнал завершения - делаем шатдаун с таймаутом
		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown:", "error", err)
		}
	case err := <-serverError:
		// еслм пришла ошибка от сервера - фаталим
		logger.Error("Server error:", "error", err)
	}

	logger.Info("Server exited")

}
