package authhandlers

import (
	"context"
	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// регистрация
func (a *AuthHandler) SignUp(c echo.Context) error {
	// получаем пользователя с фронта
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		a.logger.Error("Error in SignUp:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// валидируем входные данные
	if req.Name == "" || req.Email == "" || req.Password == "" {
		a.logger.Error("Error in SignUp:", "error", "name, email and password are required")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "name, email and password are required"})
	}

	// устанавливаем базовую роль
	req.Role = domain.RoleBasic

	// отправляем пользователя дальше в сервис
	_, err := a.useCase.CreateUser(ctx, &domain.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		a.logger.Error("Error in SignUp:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// сразу делаем логин
	response, err := a.useCase.Login(ctx, &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		a.logger.Error("Error in SignUp:", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}

	// устанавливаем аксес токен в куку
	a.setAccessTokenCookie(c, response.AccessToken)

	a.logger.Info("User signed up successful", "email", req.Email)

	return c.JSON(http.StatusCreated, models.SuccessResponse{Message: "sign up successful"})
}

// логин
func (a *AuthHandler) SignIn(c echo.Context) error {
	// получаем пользователя с фронта
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		a.logger.Error("Error in SignIn:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// отправляем дальше в сервисы
	response, err := a.useCase.Login(ctx, &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		a.logger.Error("Error in SignIn:", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}

	// устанавливаем аксес токен в куку
	a.setAccessTokenCookie(c, response.AccessToken)

	a.logger.Info("User signed in successful", "email", req.Email)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "sign in successful"})
}

// обновление токенов
func (a *AuthHandler) RefreshToken(c echo.Context) error {
	// получаем user_id из контекста
	userID, ok := c.Get("user_id").(string)
	if !ok {
		a.logger.Error("Error in RefreshToken:", "error", "unauthorized")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// обновляем токены
	response, err := a.useCase.RefreshTokens(ctx, userID)
	if err != nil {
		a.logger.Error("Error in RefreshToken:", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}

	// устанавливаем новый аксес токен в куку
	a.setAccessTokenCookie(c, response.AccessToken)

	a.logger.Info("User refresh tokens successful", "user_id", userID)

	return c.JSON(http.StatusCreated, models.SuccessResponse{Message: "tokens refreshed successful"})
}

func (a *AuthHandler) SignOut(c echo.Context) error {
	// получаем user_id из контекста
	userID, ok := c.Get("user_id").(string)
	if !ok {
		a.logger.Error("Error in SignOut:", "error", "unauthoried")
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// удаляем рефреш токен
	err := a.useCase.SignOut(ctx, userID)
	if err != nil {
		a.logger.Error("Error in SignOut:", "error", err)
	}

	// очищаем куки
	a.clearAccessTokenCookie(c)

	a.logger.Info("User signed out successful", "user_id", userID)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "successfully logged out"})
}

// установка куки
func (a *AuthHandler) setAccessTokenCookie(c echo.Context, accessToken string) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900,
	}

	c.SetCookie(cookie)
}

// удаление куки
func (a *AuthHandler) clearAccessTokenCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	c.SetCookie(cookie)
}
