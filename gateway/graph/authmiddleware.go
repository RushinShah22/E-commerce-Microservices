package graph

import (
	"errors"
	"fmt"
	"log"
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
	token, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(token string, id string) error {
	key := os.Getenv("JWT_KEY")

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		log.Panic(err)
		return err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid || id != claims["id"] {
		return errors.New("Invalid JWT.")
	}
	return nil
}
