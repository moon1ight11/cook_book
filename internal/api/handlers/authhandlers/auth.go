package authhandlers

import (
	"cook_book/internal/api/jwt"
	"cook_book/internal/services"
	"cook_book/pkg/logger"
)

// структура хэндлера
type AuthHandler struct {
	useCase services.UserUseCaseInterface
	logger  logger.Logger
	jwt     jwt.JWTInterface
}

// конструктор
func NewAuthHandler(useCase services.UserUseCaseInterface, logger logger.Logger, jwt jwt.JWTInterface) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
		logger:  logger,
		jwt:     jwt,
	}
}
