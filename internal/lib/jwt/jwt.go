package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kuromii5/sso-auth/internal/models"
)

func NewJWT(user models.User, expires time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(expires).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
