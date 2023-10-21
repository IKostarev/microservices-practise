package jwtutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTUtil struct {
	SecretKey       string        `envconfig:"SECRET_KEY" required:"true" default:"superSecretKey"`
	AccessTokenExp  time.Duration `envconfig:"ACCESS_TOKEN_EXP" required:"true" default:"15m"`
	RefreshTokenExp time.Duration `envconfig:"REFRESH_TOKEN_EXP" required:"true" default:"24h"`
}

func (ju *JWTUtil) GenerateAccessToken(userID int) (string, error) {
	return ju.generateToken(userID, ju.SecretKey, ju.AccessTokenExp)
}

func (ju *JWTUtil) GenerateRefreshToken(userID int) (string, error) {
	return ju.generateToken(userID, ju.SecretKey, ju.RefreshTokenExp)
}

func (ju *JWTUtil) VerifyToken(token string) (int, error) {
	return ju.verifyToken(token)
}

func (ju *JWTUtil) generateToken(userID int, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ju *JWTUtil) verifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return ju.SecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	exp := claims["exp"].(int64)
	if time.Unix(exp, 0).Before(time.Now()) {
		return 0, errors.New("token has expired")
	}

	userID := claims["user_id"].(int)

	return userID, nil
}
