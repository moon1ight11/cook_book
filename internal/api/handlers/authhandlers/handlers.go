package authhandlers

import (
	"context"
	"net/http"
	"time"

	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"cook_book/internal/metrics"

	"github.com/labstack/echo/v4"
)

func (a *AuthHandler) SignUp(c echo.Context) error {
	a.metrics.RecordOperation("SignUp")

	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		a.metrics.RecordError(string(metrics.ErrBadRequest), "SignUp")
		a.logger.Error("SignUp: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if req.Name == "" || req.Email == "" || req.Password == "" {
		a.metrics.RecordError(string(metrics.ErrBadRequest), "SignUp")
		a.logger.Error("SignUp: validation failed", "error", "name, email and password are required")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "name, email and password are required"})
	}

	req.Role = domain.RoleBasic

	_, err := a.useCase.CreateUser(ctx, &domain.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		a.metrics.RecordError(string(metrics.ErrInternal), "SignUp")
		a.logger.Error("SignUp: create user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	response, err := a.useCase.Login(ctx, &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		a.metrics.RecordError(string(metrics.ErrForbidden), "SignUp")
		a.logger.Error("SignUp: auto login", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	a.setAccessTokenCookie(c, response.AccessToken)

	a.logger.Info("User signed up", "email", req.Email)

	return c.JSON(http.StatusCreated, models.RefreshTokenResponse{
		RefreshToken: response.RefreshToken,
	})
}

func (a *AuthHandler) SignIn(c echo.Context) error {
	a.metrics.RecordOperation("SignIn")

	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		a.metrics.RecordError(string(metrics.ErrBadRequest), "SignIn")
		a.logger.Error("SignIn: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	response, err := a.useCase.Login(ctx, &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		a.metrics.RecordError(string(metrics.ErrForbidden), "SignIn")
		a.logger.Error("SignIn: login", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	a.setAccessTokenCookie(c, response.AccessToken)

	a.logger.Info("User signed in", "email", req.Email)

	return c.JSON(http.StatusOK, models.RefreshTokenResponse{
		RefreshToken: response.RefreshToken,
	})
}

func (a *AuthHandler) RefreshToken(c echo.Context) error {
	a.metrics.RecordOperation("RefreshToken")

	var req models.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		a.metrics.RecordError(string(metrics.ErrBadRequest), "RefreshToken")
		a.logger.Error("RefreshToken: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	response, err := a.useCase.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		a.metrics.RecordError(string(metrics.ErrForbidden), "RefreshToken")
		a.logger.Error("RefreshToken: refresh tokens", "error", err)
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	a.setAccessTokenCookie(c, response.AccessToken)

	a.logger.Info("Tokens refreshed", "user_id", response.User.Id.Hex())

	return c.JSON(http.StatusOK, models.RefreshTokenResponse{
		RefreshToken: response.RefreshToken,
	})
}

func (a *AuthHandler) SignOut(c echo.Context) error {
	a.metrics.RecordOperation("SignOut")

	userID, ok := c.Get("user_id").(string)
	if !ok {
		a.metrics.RecordError(string(metrics.ErrForbidden), "SignOut")
		a.logger.Error("SignOut: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err := a.useCase.SignOut(ctx, userID); err != nil {
		a.metrics.RecordError(string(metrics.ErrInternal), "SignOut")
		a.logger.Error("SignOut: sign out", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	a.clearAccessTokenCookie(c)

	a.logger.Info("User signed out", "user_id", userID)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "successfully logged out"})
}

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
