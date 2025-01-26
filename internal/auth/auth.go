package auth

import (
	"log"
	"time"

	"github.com/c-santos/go-auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func GenerateToken(data map[string]string) (Tokens, error) {
	key := config.LoadConfig().JWTSecret

	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"jid":  uuid.New(),
		"sub":  uuid.New(),
	}

	refreshClaims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(60 * time.Minute).Unix(),
		"jid":  uuid.New(),
		"sub":  uuid.New(),
	}

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	)

	signedAccessToken, err := accessToken.SignedString([]byte(key))
	if err != nil {
		log.Printf("[auth.go] Failed to sign access token. %s", err)
		return Tokens{}, err
	}
	log.Printf("[auth.go] Generated access token: %s", signedAccessToken)

	signedRefreshToken, err := refreshToken.SignedString([]byte(key))
	if err != nil {
		log.Printf("[auth.go] Failed to refresh sign token. %s", err)
		return Tokens{}, err
	}
	log.Printf("[auth.go] Generated refresh token: %s", signedRefreshToken)

	return Tokens{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken}, nil
}

func VerifyToken(tokens Tokens) (jwt.MapClaims, error) {
	key := config.LoadConfig().JWTSecret

	// If access token is valid, return claim 
	// If access token invalid, verify refresh token 
	// If refresh token valid, generate new access token and return along with refresh token
	// If refresh token invalid, return 401

	token, err := jwt.Parse(tokens.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err == nil || token.Valid {
		return token.Claims.(jwt.MapClaims), nil
	}

	return token.Claims.(jwt.MapClaims), nil
}
