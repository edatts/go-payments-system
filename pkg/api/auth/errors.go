package auth

import (
	"fmt"
)

func ErrUserExists(email string) error {
	return fmt.Errorf("user with email (%s) already exists", email)
}

func ErrFailedValidation(errs error) error {
	return fmt.Errorf("failed validation: %w", errs)
}
