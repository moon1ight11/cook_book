package recipehandlers

import (
	"cook_book/internal/metrics"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

type RecipeHandler struct {
	useCase services.RecipeUseCaseInterface
	logger  logger.Logger
	metrics metrics.MetricsInterface
}

func NewRecipeHandler(
	useCase services.RecipeUseCaseInterface,
	logger logger.Logger,
	metrics metrics.MetricsInterface,
) *RecipeHandler {
	return &RecipeHandler{
		useCase: useCase,
		logger:  logger,
		metrics: metrics,
	}
}
