package main

import (
	"context"
	"cook_book/internal/api"
	"cook_book/internal/api/handlers/authhandlers"
	"cook_book/internal/api/handlers/recipehandlers"
	"cook_book/internal/api/handlers/userhandlers"
	"cook_book/internal/api/jwt"
	"cook_book/internal/config"
	"cook_book/internal/services/recipe"
	"cook_book/internal/services/user"
	"cook_book/internal/storage"
	inmemcache "cook_book/internal/storage/cache/in_mem_cache"
	"cook_book/internal/storage/repo"
	"cook_book/internal/storage/repo/recipes"
	"cook_book/internal/storage/repo/users"
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
		log.Fatal("Failed to init logger: %w", err)
	}

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal("Failed to connect DB", "error:", err)
	}

	err = storage.MigrateAdminUser(db.DB, cfg)
	if err != nil {
		log.Fatal("Failed to create migrations", "error:", err)
	}

	// создание индексов в бд
	recipesColl := db.DB.Collection("recipes")
	usersColl := db.DB.Collection("users")
	indexManager := repo.NewIndexManager(recipesColl, usersColl)

	if err := indexManager.CreateAllIndexes(context.Background()); err != nil {
		logger.Warn("Warning: failed to create indexes:", "error", err)
	}

	// инициализация jwt
	jwtService := jwt.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	// инициализация кэша
	recipeCache := inmemcache.NewRecipeCache(cfg.Cache.TTL, cfg.Cache.CleanUpInterval, logger)
	defer recipeCache.Stop()

	// инициализация зависимостей
	recipeRepo := recipes.NewRecipeMongoRepo(db, logger)
	userRepo := users.NewUserMongoRepo(db, logger)

	recipeUseCase := recipe.NewRecipeUseCase(recipeRepo, recipeCache)
	userUseCase := user.NewUserUseCase(userRepo, jwtService)

	recipeHandler := recipehandlers.NewRecipeHandler(recipeUseCase, logger)
	authHandler := authhandlers.NewAuthHandler(userUseCase, logger, jwtService)
	userHandler := userhandlers.NewUserHandler(userUseCase, logger)

	// инициализируем и запускаем клинер
	cleanupService := jwt.NewCleanupService(userUseCase, logger)
	ctx := context.Background()
	cleanupService.StartCleanup(ctx)

	// инициализация роутера
	router := api.NewRouter(recipeHandler, authHandler, userHandler)

	// инициализация роутов
	router.Init(logger, jwtService)

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
		logger.Info("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// завершаем работу сервера
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown:", "error", err)
		}

		logger.Info("Disconnecting from database")

		// завершаем работу с БД
		if err := db.Close(ctx); err != nil {
			logger.Error("Failed to close database connection:", "error", err)
		} else {
			logger.Info("Database connection closed successfully")
		}

		// останавливаем очиститель токенов
		cleanupService.Stop()
	case err := <-serverError:
		// еслм пришла ошибка от сервера - фаталим
		logger.Error("Server error:", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
