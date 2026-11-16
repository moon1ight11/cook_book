package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                  string             `bson:"name" json:"name"`
	Email                 string             `bson:"email" json:"email"`
	HashPass              string             `bson:"password" json:"-"`
	Role                  Role               `bson:"role" json:"role"`
	RefreshToken          string             `bson:"refresh_token" json:"-"`
	RefreshTokenExpiresAt *time.Time         `bson:"refresh_token_expires_at" json:"-"`
	CreatedAt             time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at" json:"updated_at"`
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
