package jwt

type JWTInterface interface {
	GenerateTokens(userID, email, role string) (*TokenDetails, error)
	ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error)
	ValidateAccessToken(tokenString string) (*AccessTokenClaims, error)
}
