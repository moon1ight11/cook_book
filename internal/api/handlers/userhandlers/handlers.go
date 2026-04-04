package userhandlers

import (
	"context"
	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// получение пользователя
func (u *UserHandler) GetUser(c echo.Context) error {
	// получаем id пользователя из контекста
	id, ok := c.Get("user_id").(string)
	if !ok {
		u.logger.Error("Error in GetUser:", "error", "unauthoried")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthoried"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// идем за пользователем в сервисы
	user, err := u.useCase.GetUserByID(ctx, id)
	if err != nil {
		u.logger.Error("Error in GetUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User get successful:", "user_id", id)

	return c.JSON(http.StatusOK, user)
}

// обновление пользователя
func (u *UserHandler) UpdateUser(c echo.Context) error {
	// получаем id пользователя из контекста
	id, ok := c.Get("user_id").(string)
	if !ok {
		u.logger.Error("Error in UpdateUser:", "error", "unauthoried")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "unauthoried"})
	}

	// получаем обновленного пользователя с фронта
	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// идем в сервисы обновлять пользователя
	err := u.useCase.UpdateUser(ctx, id, req.Name, req.Email, req.Password)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User updated successful:", "user_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user updated successful"})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c echo.Context) error {
	// получаем id пользователя из контекста
	id, ok := c.Get("user_id").(string)
	if !ok {
		u.logger.Error("Error in DeleteUser:", "error", "unauthoried")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// идем в сервисы удалять пользователя
	err := u.useCase.DeleteUser(ctx, id)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User deleted successful:", "user_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successful"})
}

// ADMIN: обновление роли пользователя
func (u *UserHandler) UpdateUserRole(c echo.Context) error {
	// получаем id из параметров
	id := c.Param("user_id")

	// получаем новую роль с фронта
	var newUserRole struct {
		Role domain.Role `json:"role"`
	}
	if err := c.Bind(&newUserRole); err != nil {
		u.logger.Error("Error in UpdateUserRole:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// обновляем роль
	err := u.useCase.UpdateUserRole(ctx, newUserRole.Role, id)
	if err != nil {
		u.logger.Error("Error in UpdateUserRole:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User role updated successful:", "user", id, "role", newUserRole.Role)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user role updated successful"})
}

// ADMIN: удаление пользователя
func (u *UserHandler) DeleteAnyUser(c echo.Context) error {
	// получаем id из параметров
	id := c.Param("user_id")

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// удаляем пользователя
	err := u.useCase.DeleteAnyUser(ctx, id)
	if err != nil {
		u.logger.Error("Error in DeleteAnyUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User deleted successful:", "user", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successful"})
}

// ADMIN: создание админа
func (u *UserHandler) CreateAdmin(c echo.Context) error {
	// получаем пользователя с фронта
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		u.logger.Error("Error in CreateAdmin:", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// валидируем входные данные
	if req.Name == "" || req.Email == "" || req.Password == "" {
		u.logger.Error("Error in SignUp:", "error", "name, email and password are required")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "name, email and password are required"})
	}

	// устанавливаем роль
	req.Role = domain.RoleAdmin

	// отправляем пользователя дальше в сервис
	_, err := u.useCase.CreateUser(ctx, &domain.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		u.logger.Error("Error in CreateAdmin:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("Admin created successful", "email", req.Email)

	return c.JSON(http.StatusCreated, models.SuccessResponse{Message: "admin created successful"})
}
