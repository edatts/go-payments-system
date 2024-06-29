package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

const userKey = "userId"

// const userKey = "username"

func GetUserIdFromContext(ctx context.Context) (int32, bool) {
	val, ok := ctx.Value(userKey).(int32)
	return val, ok
}

// func GetUsernameFromContext(ctx context.Context) (string, bool) {
// 	val, ok := ctx.Value(userKey).(string)
// 	return val, ok
// }

func CreateJWT(secret []byte, userId int32) (string, error) {
	var expiration = time.Second * time.Duration(config.Envs.JWT_EXPIRATION_SECONDS)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed signing token: %w", err)
	}

	return tokenString, nil
}
