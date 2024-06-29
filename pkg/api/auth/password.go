package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("failed generating bcrypt hash: %w", err)
	}

	return string(hash), nil
}

func VerifyPassword(plaintext, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
}
