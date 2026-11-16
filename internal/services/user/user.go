package user

import (
	"cook_book/internal/api/jwt"
	"cook_book/internal/storage/repo"
)

type userUseCase struct {
	repo       repo.UserRepositoryInterface
	jwtService jwt.JWTInterface
}

func NewUserUseCase(repo repo.UserRepositoryInterface, jwt jwt.JWTInterface) *userUseCase {
	return &userUseCase{
		repo:       repo,
		jwtService: jwt,
	}
}
