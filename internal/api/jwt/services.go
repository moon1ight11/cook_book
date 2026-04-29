package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// создание токенов
func (s *JWTService) GenerateTokens(userID, email, role string) (*TokenDetails, error) {
	// обозначаем срок годности токенов
	td := &TokenDetails{
		AtExpires: time.Now().Add(s.accessExpiry).Unix(),
		RtExpires: time.Now().Add(s.refreshExpiry).Unix(),
	}

	// создание основного токена
	accessClaims := AccessTokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	var err error
	td.AccessToken, err = accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, fmt.Errorf("error in Generateokens: %w", err)
	}

	// создание рефреш-токена
	refreshClaims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return nil, fmt.Errorf("error in Generateokens: %w", err)
	}

	return td, nil
}

// декодировка и валидация основного токена
func (s *JWTService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in ValidateAccessToken: %w", err)
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("error in ValidateAccessToken: %w", err)
}

// декодировка и валидация рефреш токена
func (s *JWTService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in ValidateRefreshToken: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("error in ValidateRefreshToken: %w", err)
}
