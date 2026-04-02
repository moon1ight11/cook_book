package authhandlers

import (
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

// структура хэндлера
type AuthHandler struct {
	useCase services.UserUseCaseInterface
	logger  logger.Logger
}

// конструктор
func NewAuthHandler(useCase services.UserUseCaseInterface, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
		logger:  logger,
	}
}
