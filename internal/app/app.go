package app

import (
	"context"
	"fmt"
	"net/http"

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
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	m := metrics.NewMetrics()

	db, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, err
	}

	if err = storage.MigrateAdminUser(db.DB, cfg); err != nil {
		log.Error("Failed to run admin migration", "error", err)
		return nil, err
	}

	recipesColl := db.DB.Collection("recipes")
	usersColl := db.DB.Collection("users")
	indexManager := repo.NewIndexManager(recipesColl, usersColl)

	if err := indexManager.CreateAllIndexes(ctx); err != nil {
		log.Warn("Failed to create indexes", "error", err)
	}

	jwtService := jwt.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	recipeCache := inmemcache.NewRecipeCache(cfg.Cache.TTL, cfg.Cache.CleanUpInterval, log)

	recipeRepo := recipes.NewRecipeMongoRepo(db, log)
	userRepo := users.NewUserMongoRepo(db, log)

	recipeUseCase := recipe.NewRecipeUseCase(recipeRepo, recipeCache)
	userUseCase := user.NewUserUseCase(userRepo, jwtService)

	recipeHandler := recipehandlers.NewRecipeHandler(recipeUseCase, log, m)
	authHandler := authhandlers.NewAuthHandler(userUseCase, log, jwtService, m)
	userHandler := userhandlers.NewUserHandler(userUseCase, log, m)

	cleanupService := jwt.NewCleanupService(userUseCase, log)

	router := api.NewRouter(recipeHandler, authHandler, userHandler)
	router.Init(log, jwtService, m)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router.GetEngine(),
	}

	return &Dependencies{
		Config:         cfg,
		Logger:         log,
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

func (c *Dependencies) StartCleanup(ctx context.Context) {
	c.CleanupService.StartCleanup(ctx)
}

func (c *Dependencies) StopCleanup() {
	c.CleanupService.Stop()
}

func (c *Dependencies) StopCache() {
	c.RecipeCache.Stop()
}

func (c *Dependencies) CloseDB(ctx context.Context) error {
	return c.DB.Close(ctx)
}
