package token

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string = "3d86f180105efab801fc9178d569dbf1"

type TokenRepository interface {
	StoreRefreshToken(c context.Context, userID uint, token string) error
	RevokeToken(c context.Context, token string) error
	RevokeTokenByUserID(c context.Context, userID uint) error
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
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
	access, err := ts.generateJWT(userID, 60*time.Minute)
	if err != nil {
		return "", "", err
	}

	refresh, err := ts.generateJWT(userID, 24*7*time.Hour)
	if err != nil {
		return "", "", err
	}

	err = ts.tokenRepository.StoreRefreshToken(c, userID, refresh)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (ts *TokenService) generateJWT(userID uint, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (ts *TokenService) IsTokenValid(token string) bool {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return false
	}

	return t.Valid
}

func (ts *TokenService) DeleteToken(c context.Context, token string) error {
	return ts.tokenRepository.RevokeToken(c, token)
}

func (ts *TokenService) DeleteAllSessions(c context.Context, userID uint) error {
	return ts.tokenRepository.RevokeTokenByUserID(c, userID)
}
