package userhandlers

import (
	"context"
	"net/http"
	"time"

	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"cook_book/internal/metrics"

	"github.com/labstack/echo/v4"
)

func (u *UserHandler) GetUser(c echo.Context) error {
	u.metrics.RecordOperation("GetUser")

	id, ok := c.Get("user_id").(string)
	if !ok {
		u.metrics.RecordError(string(metrics.ErrForbidden), "GetUser")
		u.logger.Error("GetUser: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	user, err := u.useCase.GetUserByID(ctx, id)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "GetUser")
		u.logger.Error("GetUser: get user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User retrieved", "user_id", id)

	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) UpdateUser(c echo.Context) error {
	u.metrics.RecordOperation("UpdateUser")

	id, ok := c.Get("user_id").(string)
	if !ok {
		u.metrics.RecordError(string(metrics.ErrForbidden), "UpdateUser")
		u.logger.Error("UpdateUser: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateUser")
		u.logger.Error("UpdateUser: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err := u.useCase.UpdateUser(ctx, id, req.Name, req.Email, req.Password); err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "UpdateUser")
		u.logger.Error("UpdateUser: update user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User updated", "user_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user updated successful"})
}

func (u *UserHandler) DeleteUser(c echo.Context) error {
	u.metrics.RecordOperation("DeleteUser")

	id, ok := c.Get("user_id").(string)
	if !ok {
		u.metrics.RecordError(string(metrics.ErrForbidden), "DeleteUser")
		u.logger.Error("DeleteUser: missing user_id")
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err := u.useCase.DeleteUser(ctx, id); err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "DeleteUser")
		u.logger.Error("DeleteUser: delete user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User deleted", "user_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successful"})
}

func (u *UserHandler) UpdateUserRole(c echo.Context) error {
	u.metrics.RecordOperation("UpdateUserRole_ADMIN")

	id := c.Param("user_id")

	var newUserRole struct {
		Role domain.Role `json:"role"`
	}
	if err := c.Bind(&newUserRole); err != nil {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "UpdateUserRole_ADMIN")
		u.logger.Error("UpdateUserRole: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err := u.useCase.UpdateUserRole(ctx, newUserRole.Role, id); err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "UpdateUserRole_ADMIN")
		u.logger.Error("UpdateUserRole: update role", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User role updated", "user_id", id, "role", newUserRole.Role)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user role updated successful"})
}

func (u *UserHandler) DeleteAnyUser(c echo.Context) error {
	u.metrics.RecordOperation("DeleteUser_ADMIN")

	id := c.Param("user_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if err := u.useCase.DeleteAnyUser(ctx, id); err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "DeleteUser_ADMIN")
		u.logger.Error("DeleteAnyUser: delete user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("User deleted by admin", "user_id", id)

	return c.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successful"})
}

func (u *UserHandler) CreateAdmin(c echo.Context) error {
	u.metrics.RecordOperation("CreateAdmin_ADMIN")

	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "CreateAdmin_ADMIN")
		u.logger.Error("CreateAdmin: bind request", "error", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "bad request"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	if req.Name == "" || req.Email == "" || req.Password == "" {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "CreateAdmin_ADMIN")
		u.logger.Error("CreateAdmin: validation failed", "error", "name, email and password are required")
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "name, email and password are required"})
	}

	req.Role = domain.RoleAdmin

	_, err := u.useCase.CreateUser(ctx, &domain.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "CreateAdmin_ADMIN")
		u.logger.Error("CreateAdmin: create user", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
	}

	u.logger.Info("Admin created", "email", req.Email)

	return c.JSON(http.StatusCreated, models.SuccessResponse{Message: "admin created successful"})
}
