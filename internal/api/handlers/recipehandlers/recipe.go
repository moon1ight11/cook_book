package recipehandlers

import (
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

// структура хэндлера
type RecipeHandler struct {
	useCase services.RecipeUseCaseInterface
	logger  logger.Logger
}

// конструктор
func NewRecipeHandler(useCase services.RecipeUseCaseInterface, logger logger.Logger) *RecipeHandler {
	return &RecipeHandler{
		useCase: useCase,
		logger:  logger,
	}
}
