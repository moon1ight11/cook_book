package userhandlers

import (
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

// структура хэндлера
type UserHandler struct {
	useCase services.UserUseCaseInterface
	logger  logger.Logger
}

// конструктор
func NewUserHandler(useCase services.UserUseCaseInterface, logger logger.Logger) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		logger:  logger,
	}
}
