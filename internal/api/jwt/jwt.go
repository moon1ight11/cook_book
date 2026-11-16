package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTService struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTService(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}
