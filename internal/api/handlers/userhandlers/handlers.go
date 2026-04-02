package userhandlers

import (
	"context"
	"cook_book/internal/api/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// получение пользователя
func (u *UserHandler) GetUser(c echo.Context) error {
	// получаем id пользователя из контекста
	id, ok := c.Get("user_id").(string)
	if !ok {
		u.logger.Error("Error in GetUser:", "error", "unauthoried")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "unauthoried"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// идем за пользователем в сервисы
	user, err := u.useCase.GetUserByID(ctx, id)
	if err != nil {
		u.logger.Error("Error in GetUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// идем в сервисы обновлять пользователя
	err := u.useCase.UpdateUser(ctx, id, req.Name, req.Email, req.Password)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "unauthoried"})
	}

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	// идем в сервисы удалять пользователя
	err := u.useCase.DeleteUser(ctx, id)
	if err != nil {
		u.logger.Error("Error in UpdateUser:", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	u.logger.Info("User deleted successful:", "user_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successful"})
}
