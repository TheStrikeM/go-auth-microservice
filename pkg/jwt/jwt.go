package jwt

import (
	"github.com/golang-jwt/jwt"
	"microauth/internal/modules/auth/models/dto"
	"os"
)

func GenerateToken(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	secretKey := getSecretKey()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getSecretKey() string {
	if os.Getenv("JWT_SECRET_KEY") == "" {
		return "secret_key"
	}
	return os.Getenv("JWT_SECRET_KEY")
}

func GetPayload(token string) (*dto.UserDTO, error) {
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(getSecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	claims["sub"]
}
