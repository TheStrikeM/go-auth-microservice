package jwt

import (
	"github.com/golang-jwt/jwt"
	"os"
)

func GenerateToken(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if os.Getenv("JWT_SECRET_KEY") == "" {
		secretKey = "secret_key"
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
