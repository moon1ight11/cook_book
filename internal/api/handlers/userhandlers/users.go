package userhandlers

import (
	"cook_book/internal/metrics"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

type UserHandler struct {
	useCase services.UserUseCaseInterface
	logger  logger.Logger
	metrics metrics.MetricsInterface
}

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
