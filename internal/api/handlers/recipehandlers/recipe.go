package recipehandlers

import (
	"cook_book/internal/metrics"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

// структура хэндлера
type RecipeHandler struct {
	useCase services.RecipeUseCaseInterface
	logger  logger.Logger
	metrics *metrics.Metrics
}

// конструктор
func NewRecipeHandler(useCase services.RecipeUseCaseInterface, logger logger.Logger, metrics *metrics.Metrics) *RecipeHandler {
	return &RecipeHandler{
		useCase: useCase,
		logger:  logger,
		metrics: metrics,
	}
}
