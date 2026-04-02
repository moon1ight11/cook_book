package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	HashPass     string             `bson:"password"`
	Role         Role               `bson:"role"`
	RefreshToken string             `bson:"refresh_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

// модель для создания пользователя
type CreateUserRequest struct {
	Name     string `bson:"name" binding:"required"`
	Email    string `bson:"email" binding:"required,email"`
	Password string `bson:"password" binding:"required"`
	Role     Role   `bson:"role"`
}

// модель для запроса на вход
type LoginRequest struct {
	Email    string `bson:"email" binding:"required,email"`
	Password string `bson:"password" binding:"required"`
}

// модель для ответа после входа
type LoginResponse struct {
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
	User         *User  `bson:"user"`
}

type Role string

const (
	RoleBasic Role = "Basic"
	RoleAdmin Role = "Admin"
)

// ошибки юзеров
var (
	ErrEmptyName      = NewUserError("user name is required")
	ErrEmptyEmail     = NewUserError("user email is required")
	ErrEmptyPass      = NewUserError("user password is required")
	ErrNoRole         = NewUserError("user role is required")
	ErrNoRefreshToken = NewUserError("user refresh token is required")
	ErrUserNotFound   = NewUserError("user not found")
)

// структура ошибки юзеров
type UserError struct {
	Message string
}

func (e *RecipeError) UserError() string {
	return e.Message
}

// конструктор ошибки юзеров
func NewUserError(msg string) *RecipeError {
	return &RecipeError{Message: msg}
}
