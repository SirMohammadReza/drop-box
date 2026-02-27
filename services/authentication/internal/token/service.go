package token

import (
	jwtauth "authentication/internal/platform/jwt_auth"
	"context"
	"fmt"
)

type TokenType int

func (t TokenType) String() string {
	switch t {
	case 1:
		return "Access Token Type"
	case 2:
		return "Refresh Token Type"
	default:
		return "Invalid Token Type"
	}
}

const (
	TokenAccessType  = TokenType(1)
	TokenRefreshType = TokenType(2)
)

type TokenRepository interface {
	StoreRefreshToken(c context.Context, token *Token) error
	RevokeToken(c context.Context, token string) error
}

type TokenService struct {
	tokenRepository TokenRepository
}

func NewTokenService(tr TokenRepository) *TokenService {
	return &TokenService{
		tokenRepository: tr,
	}
}

func (ts *TokenService) GenerateTokenPair(c context.Context, userID uint) (string, string, error) {
	token, refreshToken, err := jwtauth.GenerateTokenPair(userID)
	if err != nil {
		return "", "", fmt.Errorf("could not generate tokens: %w", err)
	}

	err = ts.tokenRepository.StoreRefreshToken(c, &Token{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return "", "", fmt.Errorf("could not store resfresh token: %w", err)
	}

	return token, refreshToken, nil
}

func (ts *TokenService) ValidateToken(c context.Context, token string, tokenType TokenType) (*uint, error) {
	if tokenType == TokenAccessType {
		claims, err := jwtauth.ValidateAccessToken(token)
		if err != nil {
			return nil, fmt.Errorf("could not validate access token: %w", err)
		}

		return &claims.UserID, nil
	}

	if tokenType == TokenRefreshType {
		claims, err := jwtauth.ValidateRefreshToken(token)
		if err != nil {
			return nil, fmt.Errorf("could not validate refresh token: %w", err)
		}

		return &claims.UserID, nil
	}

	return nil, fmt.Errorf("could not validate token, invalid token type: %s", tokenType)
}

func (ts *TokenService) DeleteToken(c context.Context, token string) error {
	return ts.tokenRepository.RevokeToken(c, token)
}
