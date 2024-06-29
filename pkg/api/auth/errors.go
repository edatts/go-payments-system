package auth

import (
	"errors"
	"fmt"
)

var (
	ErrWrongPassword = errors.New("incorrect password for user")
)

func ErrEmailExists(email string) error {
	return fmt.Errorf("user with email (%s) already exists", email)
}

func ErrEmailNotExists(email string) error {
	return fmt.Errorf("no user found with email address (%s)", email)
}

func ErrUsernameExists(username string) error {
	return fmt.Errorf("user with username (%s) already exists", username)
}

func ErrUsernameNotExists(username string) error {
	return fmt.Errorf("no user found with username (%s)", username)
}

func ErrFailedValidation(errs error) error {
	return fmt.Errorf("failed validation: %w", errs)
}
