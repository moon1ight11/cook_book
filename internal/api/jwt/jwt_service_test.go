package jwt

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJWTService() *JWTService {
	return &JWTService{
		accessSecret:  "test_access_secret_key",
		refreshSecret: "test_refresh_secret_key",
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
	}
}

// успешное создание токенов
func TestGenerateTokens_Success(t *testing.T) {
	s := setupJWTService()

	userID := "507f1f77bcf86cd799439011"
	email := "test@example.com"
	role := "basic"

	td, err := s.GenerateTokens(userID, email, role)

	require.NoError(t, err)
	require.NotNil(t, td)

	// токены не пустые
	assert.NotEmpty(t, td.AccessToken)
	assert.NotEmpty(t, td.RefreshToken)

	// токены разные
	assert.NotEqual(t, td.AccessToken, td.RefreshToken)

	// сроки действия
	assert.True(t, td.AtExpires > time.Now().Unix())
	assert.True(t, td.RtExpires > time.Now().Unix())
	assert.True(t, td.RtExpires > td.AtExpires)

	// access истекает через 15 минут
	expectedAccessExpiry := time.Now().Add(15 * time.Minute).Unix()
	assert.InDelta(t, expectedAccessExpiry, td.AtExpires, 0)

	// refresh истекает через 7 дней
	expectedRefreshExpiry := time.Now().Add(7 * 24 * time.Hour).Unix()
	assert.InDelta(t, expectedRefreshExpiry, td.RtExpires, 0)
}

// создание токенов для разных пользователей
func TestGenerateTokens_DifferentUsers(t *testing.T) {
	s := setupJWTService()

	td1, err := s.GenerateTokens("user1", "user1@test.com", "basic")
	require.NoError(t, err)

	td2, err := s.GenerateTokens("user2", "user2@test.com", "admin")
	require.NoError(t, err)

	// токены разные
	assert.NotEqual(t, td1.AccessToken, td2.AccessToken)
	assert.NotEqual(t, td1.RefreshToken, td2.RefreshToken)

	// клеймы первого юзера
	claims1, err := s.ValidateAccessToken(td1.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "user1", claims1.UserID)
	assert.Equal(t, "basic", claims1.Role)

	// клеймы второго юзера
	claims2, err := s.ValidateAccessToken(td2.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "user2", claims2.UserID)
	assert.Equal(t, "admin", claims2.Role)
}

// успешная валидация аксес токена
func TestValidateAccessToken_Success(t *testing.T) {
	s := setupJWTService()

	td, err := s.GenerateTokens("user789", "user@test.com", "basic")
	require.NoError(t, err)

	claims, err := s.ValidateAccessToken(td.AccessToken)

	require.NoError(t, err)
	assert.Equal(t, "user789", claims.UserID)
	assert.Equal(t, "user@test.com", claims.Email)
	assert.Equal(t, "basic", claims.Role)
}

// валидация неверного аксес токена
func TestValidateAccessToken_InvalidToken(t *testing.T) {
	s := setupJWTService()

	claims, err := s.ValidateAccessToken("invalid.token.here")
	
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// валидация аксес токена с не нашим секретом
func TestValidateAccessToken_WrongSecret(t *testing.T) {
	s := setupJWTService()

	td, err := s.GenerateTokens("user123", "test@test.com", "basic")
	require.NoError(t, err)

	s2 := &JWTService{
		accessSecret:  "WRONGSECRETKEY",
		refreshSecret: "test_refresh_secret_key",
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
	}

	claims, err := s2.ValidateAccessToken(td.AccessToken)
	
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// валидация просроченного аксес токена
func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	s := &JWTService{
		accessSecret:  "test_access_secret",
		refreshSecret: "test_refresh_secret",
		accessExpiry:  -1 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
	}

	td, err := s.GenerateTokens("user123", "test@test.com", "basic")
	require.NoError(t, err)

	claims, err := s.ValidateAccessToken(td.AccessToken)
	
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

// успешная валидация рефреш токена
func TestValidateRefreshToken_Success(t *testing.T) {
	s := setupJWTService()

	td, err := s.GenerateTokens("user", "test@test.com", "admin")
	require.NoError(t, err)

	claims, err := s.ValidateRefreshToken(td.RefreshToken)

	require.NoError(t, err)
	assert.Equal(t, "user", claims.UserID)
}

// валидация неправильного рефреша
func TestValidateRefreshToken_InvalidToken(t *testing.T) {
	s := setupJWTService()

	claims, err := s.ValidateRefreshToken("invalid.refresh.token")
	
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// валидация рефреша с не нашим секретом
func TestValidateRefreshToken_WrongSecret(t *testing.T) {
	s := setupJWTService()

	td, err := s.GenerateTokens("user4", "test@test.com", "basic")
	require.NoError(t, err)

	s2 := &JWTService{
		accessSecret:  "test_access_secret_key",
		refreshSecret: "WRONGSECRET",
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
	}

	claims, err := s2.ValidateRefreshToken(td.RefreshToken)
	
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// валидация рефреша прсроченного
func TestValidateRefreshToken_ExpiredToken(t *testing.T) {
	s := &JWTService{
		accessSecret:  "test_access_secret",
		refreshSecret: "test_refresh_secret",
		accessExpiry:  15 * time.Minute,
		refreshExpiry: -1 * time.Hour,
	}

	td, err := s.GenerateTokens("user5", "test@test.com", "basic")
	require.NoError(t, err)

	claims, err := s.ValidateRefreshToken(td.RefreshToken)
	
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

// генерация токеов с пустыми полями
func TestGenerateTokens_EmptyFields(t *testing.T) {
	s := setupJWTService()

	td, err := s.GenerateTokens("", "", "")
	require.NoError(t, err)
	require.NotNil(t, td)

	claims, err := s.ValidateAccessToken(td.AccessToken)
	require.NoError(t, err)
	assert.Empty(t, claims.UserID)
	assert.Empty(t, claims.Email)
	assert.Empty(t, claims.Role)
}
