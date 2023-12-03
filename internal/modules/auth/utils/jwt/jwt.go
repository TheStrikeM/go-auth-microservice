package jwt

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func GenerateToken(id, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(2 * time.Hour)
	claims["sub"] = id
	claims["username"] = username

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
