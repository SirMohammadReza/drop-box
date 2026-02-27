package jwtauth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	refreshSecret = []byte("B5MRdKX1V0V9ED506XlkMKWnaEEYyJmErGziWEMeEDg=")
	accessSecret  = []byte("kr1Fst3dceSEHciyFJSLSVL9A08+ea5Xu/b1RUjxr6U=")
)

type Claims struct {
	Uuid uuid.UUID `json:"uuid"`
	jwt.RegisteredClaims
}

func generateToken(uuid uuid.UUID, secret []byte, expireAt time.Duration) (string, error) {
	claims := &Claims{
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func GenerateTokenPair(uuid uuid.UUID) (string, string, error) {
	accessToken, err := generateToken(uuid, accessSecret, time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateToken(uuid, refreshSecret, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse refresh token with claims: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return accessSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token with claims: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
