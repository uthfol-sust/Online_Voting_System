package utils

import (
	"fmt"
	"pollvoting/pkg/config"
	"pollvoting/pkg/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyCustomClaims struct {
	ID   int64
	Name string
	Role string
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	ID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *models.User) (string, error) {
	key := config.LocalConfig.JWTsecret

	claims := MyCustomClaims{
		ID:   user.ID,
		Name: user.Name,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "pollVoting",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func VerifyByJWT(tokenString string) (*MyCustomClaims, error) {
	secret_key := config.LocalConfig.JWTsecret

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*MyCustomClaims)

	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}
	return claims, nil
}

func GenerateRefreshToken(user *models.User) (string, error) {
	claims := RefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "pollVoting",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(config.LocalConfig.JWTrefreshSecret))
}

func VerifyRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.LocalConfig.JWTrefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)

	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}
	return claims, nil
}
