package utils

import (
	"errors"
	"fmt"
	"strconv"
	"template/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int) (string, error) {
	return generateToken(userID, time.Duration(12)*time.Hour, "access")
}

func GenerateRefreshToken(userID int) (string, error) {
	return generateToken(userID, time.Duration(43200)*time.Minute, "refresh")
}

func generateToken(userID int, ttl time.Duration, tokenType string) (string, error) {
	now := time.Now().UTC()

	claims := CustomClaims{
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := configs.Envs.JWTSecret
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not configured")
	}

	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	secret := configs.Envs.JWTSecret
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not configured")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
