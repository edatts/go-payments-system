package jwt

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/edatts/go-payment-system/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type CtxKey string

const UserKey = CtxKey("userId")
const JWTKey = CtxKey("JWTPublicKey")

func GetUserIdFromContext(ctx context.Context) (int32, bool) {
	val, ok := ctx.Value(UserKey).(int32)
	return val, ok
}

func GetJWTPublicKeyFromContext(ctx context.Context) ([]byte, bool) {
	val, ok := ctx.Value(JWTKey).([]byte)
	return val, ok
}

// func GetUsernameFromContext(ctx context.Context) (string, bool) {
// 	val, ok := ctx.Value(userKey).(string)
// 	return val, ok
// }

// func CreateJWT(secret []byte, userId int32) (string, error) {
// 	var expiration = time.Second * time.Duration(config.Envs.JWT_EXPIRATION_SECONDS)

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userId":    userId,
// 		"expiresAt": time.Now().Add(expiration).Unix(),
// 	})

// 	tokenString, err := token.SignedString(secret)
// 	if err != nil {
// 		return "", fmt.Errorf("failed signing token: %w", err)
// 	}

// 	return tokenString, nil
// }

func CreateJWT(secret ed25519.PrivateKey, userId int32) (string, error) {
	var expiration = time.Second * time.Duration(config.Envs.JWT_EXPIRATION_SECONDS)

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"userId":    userId,
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed signing token: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, pubkey ed25519.PublicKey) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}

		return pubkey, nil
	})
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		// Get token
		tokenString := utils.GetTokenFromRequest(req)

		// Get cached jwt public key
		pubkey, ok := GetJWTPublicKeyFromContext(req.Context())
		if !ok {
			// Log it then permission denied
			log.Printf("failed getting public key from context: failed type assertion")
			utils.WriteError(rw, http.StatusInternalServerError)
			return
		}

		// Validate token
		token, err := ValidateJWT(tokenString, pubkey)
		if err != nil {
			log.Printf("failed validating token: %s", err)
			utils.WriteError(rw, http.StatusUnauthorized)
			return
		}

		// Get UserId from claims
		claims := token.Claims.(jwt.MapClaims)
		userIdStr := claims["userId"].(string)

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Printf("failed converting userId to int")
			utils.WriteError(rw, http.StatusUnauthorized)
			return
		}

		// Get user
		user, err := store.GetUserById(int32(userId))
		if err != nil {
			log.Printf("failed getting user by id: %s", err)
			utils.WriteError(rw, http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(req.Context(), UserKey, user.Id)
		req = req.WithContext(ctx)

		handlerFunc(rw, req)
	}
}
