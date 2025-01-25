package auth

import (
	"log"
	"time"

	"github.com/c-santos/go-auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(data map[string]string) (string, error) {
	key := config.LoadConfig().JWTSecret

	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		log.Printf("[auth.go] Failed to sign token. %s", err)
		return "", err
	}
	log.Printf("[auth.go] Generated token: %s", signedToken)

	return signedToken, nil
}
