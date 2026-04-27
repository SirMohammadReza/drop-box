package token

import (
	jwtauth "authentication/internal/platform/jwt_auth"
	"authentication/internal/platform/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
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

type Provider interface {
	GenerateTokenPair(c context.Context, userID uint, uuid uuid.UUID) (string, string, error)
	ValidateToken(c context.Context, token string, tokenType TokenType) (*uuid.UUID, error)
	DeleteToken(c context.Context, token string) error
	RefreshTokens(c context.Context, refreshToken string) (string, string, error)
}

type TokenRepository interface {
	StoreRefreshToken(c context.Context, token *Token) error
	RevokeToken(c context.Context, token string) error
}

type TokenService struct {
	tokenRepository TokenRepository
	logger          logger.Logger
}

func NewTokenService(tr TokenRepository, l logger.Logger) Provider {
	return &TokenService{
		tokenRepository: tr,
		logger:          l,
	}
}

func (ts *TokenService) GenerateTokenPair(c context.Context, userID uint, uuid uuid.UUID) (string, string, error) {
	token, refreshToken, err := jwtauth.GenerateTokenPair(uuid)
	if err != nil {
		ts.logger.WithField("func", "generate token pair").
			WithField("user_id", userID).
			WithField("error", err).
			Errorf("could not generate pair token: %s", err.Error())
		return "", "", fmt.Errorf("could not generate tokens: %w", err)
	}

	err = ts.tokenRepository.StoreRefreshToken(c, &Token{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		ts.logger.WithField("func", "generate token pair").
			WithField("user_id", userID).
			WithField("error", err).
			Errorf("could not store refresh token: %s", err.Error())
		return "", "", fmt.Errorf("could not store resfresh token: %w", err)
	}

	return token, refreshToken, nil
}

func (ts *TokenService) ValidateToken(c context.Context, token string, tokenType TokenType) (*uuid.UUID, error) {
	if tokenType == TokenAccessType {
		claims, err := jwtauth.ValidateAccessToken(token)
		if err != nil {
			ts.logger.WithField("func", "validate token").
				WithField("token info", []string{token, tokenType.String()}).
				WithField("error", err).
				Errorf("could not validate access token: %s", err.Error())
			return nil, fmt.Errorf("could not validate access token: %w", err)
		}

		return &claims.Uuid, nil
	}

	if tokenType == TokenRefreshType {
		claims, err := jwtauth.ValidateRefreshToken(token)
		if err != nil {
			ts.logger.WithField("func", "validate token").
				WithField("token info", []string{token, tokenType.String()}).
				WithField("error", err).
				Errorf("could not validate refresh token: %s", err.Error())
			return nil, fmt.Errorf("could not validate refresh token: %w", err)
		}

		return &claims.Uuid, nil
	}

	ts.logger.WithField("func", "validate token").Errorf("invalud token type: %s", tokenType)
	return nil, fmt.Errorf("could not validate token, invalid token type: %s", tokenType)
}

func (ts *TokenService) DeleteToken(c context.Context, token string) error {
	return ts.tokenRepository.RevokeToken(c, token)
}

func (ts *TokenService) RefreshTokens(c context.Context, refreshToken string) (string, string, error) {
	uuid, err := ts.ValidateToken(c, refreshToken, TokenRefreshType)
	if err != nil {
		ts.logger.WithField("func", "refresh token").
			WithField("error", err).Errorf("could not validate refresh token: %s", err.Error())
		return "", "", fmt.Errorf("could not validate refresh token: %w", err)
	}

	err = ts.DeleteToken(c, refreshToken)
	if err != nil {
		ts.logger.WithField("func", "refresh token").
			WithField("error", err).Errorf("could not delete refresh token: %s", err.Error())
		return "", "", fmt.Errorf("could not delete refresh token: %w", err)
	}

	acsToken, refToken, err := jwtauth.GenerateTokenPair(*uuid)
	if err != nil {
		ts.logger.WithField("func", "refresh token").
			WithField("error", err).Errorf("could not generate access and refresh token: %s", err.Error())
		return "", "", fmt.Errorf("could not generate access and refresh tokens: %w", err)
	}

	return acsToken, refToken, nil
}
