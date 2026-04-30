package app

import (
	"context"
	"cook_book/internal/api"
	"cook_book/internal/api/handlers"
	"cook_book/internal/api/handlers/authhandlers"
	"cook_book/internal/api/handlers/recipehandlers"
	"cook_book/internal/api/handlers/userhandlers"
	"cook_book/internal/api/jwt"
	"cook_book/internal/config"
	"cook_book/internal/metrics"
	"cook_book/internal/services"
	"cook_book/internal/services/recipe"
	"cook_book/internal/services/user"
	"cook_book/internal/storage"
	"cook_book/internal/storage/cache"
	inmemcache "cook_book/internal/storage/cache/in_mem_cache"
	"cook_book/internal/storage/repo"
	"cook_book/internal/storage/repo/recipes"
	"cook_book/internal/storage/repo/users"
	"cook_book/pkg/logger"
	"fmt"
	"net/http"
)

type Dependencies struct {
	Config         *config.Config
	Logger         logger.Logger
	DB             *storage.DataBase
	JWTService     jwt.JWTInterface
	RecipeCache    cache.RecipeCacheInterface
	RecipeRepo     repo.RecipeRepositoryInterface
	UsersRepo      repo.UserRepositoryInterface
	RecipeUseCase  services.RecipeUseCaseInterface
	UserUseCase    services.UserUseCaseInterface
	RecipeHandler  handlers.RecipeHandlersInterface
	AuthHandler    handlers.AuthHandlersInterface
	UserHandler    handlers.UserHandlersInterface
	Router         *api.Router
	CleanupService *jwt.CleanupService
	Server         *http.Server
	Metrics        metrics.Metrics
}

func NewDependencies(ctx context.Context) (*Dependencies, error) {
	// инициализация конфига
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// инициализация логгера
	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	metrics := metrics.NewMetrics()

	// соединение с БД
	db, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, err
	}

	// применение тн миграций
	err = storage.MigrateAdminUser(db.DB, cfg)
	if err != nil {
		logger.Error("Failed to create migrations", "error", err)
		return nil, err
	}

	// создание индексов в бд
	recipesColl := db.DB.Collection("recipes")
	usersColl := db.DB.Collection("users")
	indexManager := repo.NewIndexManager(recipesColl, usersColl)

	if err := indexManager.CreateAllIndexes(ctx); err != nil {
		logger.Warn("Warning: failed to create indexes:", "error", err)
	}

	// инициализация jwt
	jwtService := jwt.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	// инициализация кэша
	recipeCache := inmemcache.NewRecipeCache(cfg.Cache.TTL, cfg.Cache.CleanUpInterval, logger)

	// инициализация зависимостей
	recipeRepo := recipes.NewRecipeMongoRepo(db, logger)
	userRepo := users.NewUserMongoRepo(db, logger)

	recipeUseCase := recipe.NewRecipeUseCase(recipeRepo, recipeCache)
	userUseCase := user.NewUserUseCase(userRepo, jwtService)

	recipeHandler := recipehandlers.NewRecipeHandler(recipeUseCase, logger, metrics)
	authHandler := authhandlers.NewAuthHandler(userUseCase, logger, jwtService, metrics)
	userHandler := userhandlers.NewUserHandler(userUseCase, logger, metrics)

	// инициализируем клинер
	cleanupService := jwt.NewCleanupService(userUseCase, logger)

	// инициализация роутера
	router := api.NewRouter(recipeHandler, authHandler, userHandler)

	// инициализация роутов
	router.Init(logger, jwtService, metrics)

	// создаем HTTP сервер
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router.GetEngine(),
	}

	return &Dependencies{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		JWTService:     jwtService,
		RecipeCache:    recipeCache,
		RecipeRepo:     recipeRepo,
		UsersRepo:      userRepo,
		RecipeUseCase:  recipeUseCase,
		UserUseCase:    userUseCase,
		RecipeHandler:  recipeHandler,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		Router:         router,
		CleanupService: cleanupService,
		Server:         server,
	}, nil
}

// StartCleanup запускает очистку токенов
func (c *Dependencies) StartCleanup(ctx context.Context) {
	c.CleanupService.StartCleanup(ctx)
}

// StopCleanup останавливает очистку токенов
func (c *Dependencies) StopCleanup() {
	c.CleanupService.Stop()
}

// StopCache останавливает кэш
func (c *Dependencies) StopCache() {
	c.RecipeCache.Stop()
}

// CloseDB закрывает соединение с БД
func (c *Dependencies) CloseDB(ctx context.Context) error {
	return c.DB.Close(ctx)
}
