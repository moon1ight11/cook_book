package user

import (
	"context"
	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// создание пользователя
func (u *userUseCase) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// проверка, нет ли юзеров с таким имейлом
	existingUser, _ := u.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("error in CreateUser: user with email %s already exist", req.Email)
	}

	// проверяем имейл на правильность
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error in UpdateUser: %w", err)
	}
	if !matched {
		return nil, fmt.Errorf("error in UpdateUser: new email not looks like email")
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
		return nil, fmt.Errorf("error in CreateUser: %w", err)
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

	// инвалидируем старый токен
	err = u.repo.UpdateRefreshToken(ctx, user.Id, "")
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

	// хэштруем токен
	hash := sha256.Sum256([]byte(tokens.RefreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	// закидываем рефрештокен в бд
	err = u.repo.UpdateRefreshToken(ctx, user.Id, hashedToken)
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
func (u *userUseCase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
	//  валидируем токен и достаем клеймы
	claims, err := u.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	// приводим айди пользоватля к правильному виду
	objID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	// находим пользователя
	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	// проверяем что рефреш в бд совпадает с переданным
	hash := sha256.Sum256([]byte(refreshToken))
	requestHash := hex.EncodeToString(hash[:])

	if subtle.ConstantTimeCompare([]byte(user.RefreshToken), []byte(requestHash)) != 1 {
		return nil, fmt.Errorf("error in RefreshTokens: invalid token")
	}

	// инвалидируем старый токен
	err = u.repo.UpdateRefreshToken(ctx, user.Id, "")
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
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

	// хэшируем токен
	hash = sha256.Sum256([]byte(tokens.RefreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	// обновляем рефрештокен в бд
	err = u.repo.UpdateRefreshToken(ctx, user.Id, hashedToken)
	if err != nil {
		return nil, fmt.Errorf("error in RefreshTokens: %w", err)
	}

	return &domain.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

// получение пользователя по id
func (u *userUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("error in GetUserByID: %w", err)
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
		return fmt.Errorf("error in UpdateUser: %w", err)
	}

	var hashedPass string

	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return fmt.Errorf("error in UpdateUser: %w", err)
	}

	needsNewTokens := false

	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return fmt.Errorf("error in UpdateUser: new name cant be empty")
		}
	}

	finalemail := user.Email
	if email != nil && *email != "" && *email != user.Email {
		// проверяем, чтобы почта не была уже занята
		existingUser, _ := u.repo.FindByEmail(ctx, *email)
		if existingUser != nil && existingUser.Id != user.Id {
			return fmt.Errorf("email %s already taken", *email)
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

		needsNewTokens = true
		finalemail = *email
	}

	var hashedPassPtr *string
	if password != nil {
		if strings.TrimSpace(*password) == "" {
			return fmt.Errorf("error in UpdateUser: new password cant be empty")
		}
		// хэшируем пароль
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
		hashedPass = string(hashedBytes)
		hashedPassPtr = &hashedPass
		needsNewTokens = true
	}

	err = u.repo.Update(ctx, user.Id, name, email, hashedPassPtr)
	if err != nil {
		return fmt.Errorf("error in UpdateUser: %w", err)
	}

	// если изменены чувствительные поля - нужны новые токены
	if needsNewTokens {
		tokens, err := u.jwtService.GenerateTokens(
			user.Id.Hex(),
			finalemail,
			string(user.Role),
		)
		if err != nil {
			return fmt.Errorf("error in UpdateUser: %w", err)
		}

		// хэшируем токен
		hash := sha256.Sum256([]byte(tokens.RefreshToken))
		hashedToken := hex.EncodeToString(hash[:])

		err = u.repo.UpdateRefreshToken(ctx, user.Id, hashedToken)
		if err != nil {
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
	}

	return nil
}

// обновление роли пользователя
func (u *userUseCase) UpdateUserRole(ctx context.Context, newRole domain.Role, id string) error {
	// проверяем чтобы роль была правильная
	if newRole != domain.RoleAdmin && newRole != domain.RoleBasic {
		return fmt.Errorf("errror in UpdateUserRole: new role not valid: %s", string(newRole))
	}

	// приводим id к правильному типу
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error in UpdateUserRole: %w", err)
	}

	// обновляем роль
	err = u.repo.UpdateRole(ctx, string(newRole), objID)
	if err != nil {
		return fmt.Errorf("error in UpdateUserRole: %w", err)
	}

	return nil
}

// удаление пользователя
func (u *userUseCase) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error in DeleteUser: %w", err)
	}

	err = u.repo.Delete(ctx, objID)
	if err != nil {
		return fmt.Errorf("error in DeleteUser: %w", err)
	}

	return nil
}

// удаление пользователя для админа
func (u *userUseCase) DeleteAnyUser(ctx context.Context, id string) error {
	// приводим id к правильному типу
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error in DeleteAnyUser: %w", err)
	}

	// удаляем пользователя
	err = u.repo.DeleteAny(ctx, objID)
	if err != nil {
		return fmt.Errorf("error in DeleteAnyUser: %w", err)
	}

	return nil
}

// выход из профиля
func (u *userUseCase) SignOut(ctx context.Context, userID string) error {
	// получаем id из контекста
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("error in SignOut: %w", err)
	}

	// удаляем рефреш токен из БД
	err = u.repo.UpdateRefreshToken(ctx, objID, "")
	if err != nil {
		return fmt.Errorf("error in SignOut: %w", err)
	}

	return nil
}

// очистка просроченных токенов
func (u *userUseCase) CleanExpiredTokens(ctx context.Context) error {
	err := u.repo.CleanExpiredRefreshTokens(ctx)
	if err != nil {
		return fmt.Errorf("error in CleanExpiredTokens: %w", err)
	}
	return nil
}
