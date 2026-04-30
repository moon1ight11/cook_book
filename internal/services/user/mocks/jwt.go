package mocks

import (
	"cook_book/internal/api/jwt"

	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateTokens(userID, email, role string) (*jwt.TokenDetails, error) {
	args := m.Called(userID, email, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenDetails), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*jwt.RefreshTokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RefreshTokenClaims), args.Error(1)
}

func (m *MockJWTService) ValidateAccessToken(token string) (*jwt.AccessTokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.AccessTokenClaims), args.Error(1)
}
