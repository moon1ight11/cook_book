package user

import (
	"context"
	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

// создание пользователя
func (u *userUseCase) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// проверка, нет ли юзеров с таким имейлом
	existingUser, _ := u.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("error in CreateUser: user with email %s already exist", req.Email)
	}

	// хэшируем пароль
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error in CreateUser: %w", err)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		HashPass: string(hashedPass),
		Role:     req.Role,
	}

	// создаем пользователя
	err = u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// вход пользователя
func (u *userUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*models.LoginResponse, error) {
	// ищем пользователя по имейл
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error in Login: %w", err)
	}

	// сверяем пароли
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("error in Login: %w", err)
	}

	// создаем токены
	tokens, err := u.jwtService.GenerateTokens(
		user.Id.Hex(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, fmt.Errorf("error in Login: %w", err)
	}

	// закидываем рефрештокен в бд
	err = u.repo.UpdateRefreshToken(ctx, user.Id, tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error in Login: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

// обновление токенов
func (u *userUseCase) RefreshTokens(ctx context.Context, userID string) (*domain.LoginResponse, error) {
	// конвертируем ID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// находим пользователя по айди
	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	// проверяем, что есть рефреш токен
	if user.RefreshToken == "" {
		return nil, fmt.Errorf("error in RefreshTokens: no refresh token found")
	}

	// генерируем новую пару токенов
	tokens, err := u.jwtService.GenerateTokens(
		user.Id.Hex(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	// обновляем рефрештокен в бд
	err = u.repo.UpdateRefreshToken(ctx, user.Id, tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	return &domain.LoginResponse{
		AccessToken:  tokens.AccessToken,
		User:         user,
	}, nil
}

// получение пользователя по id
func (u *userUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("error in GetUserByID: %w", err)
	}

	return user, nil
}

func (u *userUseCase) UpdateUser(
	ctx context.Context,
	id string,
	name *string,
	email *string,
	password *string,
) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrUserNotFound
	}

	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return fmt.Errorf("error in UpdateUser: %w", err)
	}

	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return fmt.Errorf("error in UpdateUser: new name cant be empty")
		}
	}
	if email != nil {
		if strings.TrimSpace(*email) == "" {
			return fmt.Errorf("error in UpdateUser: new email cant be empty")
		}
		// почта должна быть похожа на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *email)
		if err != nil {
			return fmt.Errorf("error in UpdateUser: %w", err)
		}

		if !matched {
			return fmt.Errorf("error in UpdateUser: new email not looks like email")
		}
	}
	if password != nil {
		if strings.TrimSpace(*password) == "" {
			return fmt.Errorf("error in UpdateUser: new password cant be empty")
		}
	}

	err = u.repo.Update(ctx, user.Id, name, email, password)
	if err != nil {
		return fmt.Errorf("error in UpdateUser: %w", err)
	}

	return nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrUserNotFound
	}

	err = u.repo.Delete(ctx, objID)
	if err != nil {
		return fmt.Errorf("error in DeleteUser: %w", err)
	}

	return nil
}

func (u *userUseCase) SignOut(ctx context.Context, userID string) error {
	// получаем id из контекста	
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }

	// удаляем рефреш токен из БД
	err = u.repo.UpdateRefreshToken(ctx, objID, "")
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}
