package userhandlers

import (
	"cook_book/internal/metrics"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

// структура хэндлера
type UserHandler struct {
	useCase services.UserUseCaseInterface
	logger  logger.Logger
	metrics metrics.MetricsInterface
}

// конструктор
func NewUserHandler(
	useCase services.UserUseCaseInterface,
	logger logger.Logger,
	metrics metrics.MetricsInterface,
) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		logger:  logger,
		metrics: metrics,
	}
}
