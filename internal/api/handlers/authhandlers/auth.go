package authhandlers

import (
	"cook_book/internal/api/jwt"
	"cook_book/internal/metrics"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

type AuthHandler struct {
	useCase services.UserUseCaseInterface
	logger  logger.Logger
	jwt     jwt.JWTInterface
	metrics metrics.MetricsInterface
}

func NewAuthHandler(
	useCase services.UserUseCaseInterface,
	logger logger.Logger,
	jwt jwt.JWTInterface,
	metrics metrics.MetricsInterface,
) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
		logger:  logger,
		jwt:     jwt,
		metrics: metrics,
	}
}
