package graph

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// This function generates a new jwt
func GenerateToken(id string, role string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sus": []string{id, role},
		"exp": time.Now().Add(time.Hour * 24 * 2).Unix(),
		"iat": time.Now().Unix(),
	})
	key := os.Getenv("JWT_KEY")
	token, err := claims.SignedString([]byte(key))
	if err != nil {
		log.Panic(err)
		return "", err
	}
	return token, nil
}

// This function Verifies the JWT token of the user
func VerifyToken(token string) (string, string, error) {

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	user, ok := t.Claims.(jwt.MapClaims)

	id := user["sus"].([]interface{})[0].(string)
	role := user["sus"].([]interface{})[1].(string)

	if !ok || !t.Valid {
		log.Panic("incorrect verification")
		return "", "", fmt.Errorf("validation failed. Please login again.")
	}
	return id, role, nil
}
