package tokenmanager

import (
	"auth-service/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateUserToken(user *domain.User, secret []byte) (string, error ){
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return  signedToken, nil
}