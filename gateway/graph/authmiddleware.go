package graph

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 2).Unix(),
		"iat": time.Now().Unix(),
	})
	key := os.Getenv("JWT_KEY")
	token, err := claims.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(token string) error {
	key := os.Getenv("JWT_KEY")
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return interface{}(key), nil
	})

	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("Invalid JWT.")
	}
	return nil
}
