package jwtauth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	refreshSecret = []byte("B5MRdKX1V0V9ED506XlkMKWnaEEYyJmErGziWEMeEDg=")
	accessSecret  = []byte("kr1Fst3dceSEHciyFJSLSVL9A08+ea5Xu/b1RUjxr6U=")
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func generateToken(userID uint, secret []byte, expireAt time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func GenerateTokenPair(userID uint) (string, string, error) {
	accessToken, err := generateToken(userID, accessSecret, time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateToken(userID, refreshSecret, 7*24*time.Hour)
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
		return refreshSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token with claims: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
