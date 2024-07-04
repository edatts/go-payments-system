package auth

import (
	"errors"
	"fmt"
)

var (
	ErrWrongPassword     = errors.New("incorrect password for user")
	ErrEmailExists       = errors.New("user with email already exists")
	ErrEmailNotExists    = errors.New("no user found with email address")
	ErrUsernameExists    = errors.New("user with username already exists")
	ErrUsernameNotExists = errors.New("no user found with username")
	ErrFailedValidation  = errors.New("failed validation")
)

func EmailExistsError(email string) error {
	return fmt.Errorf("%w: %s", ErrEmailExists, email)
}

func EmailNotExistsError(email string) error {
	return fmt.Errorf("%w: %s", ErrEmailNotExists, email)
}

func UsernameExistsError(username string) error {
	return fmt.Errorf("%w: %s", ErrUsernameExists, username)
}

func UsernameNotExistsError(username string) error {
	return fmt.Errorf("%w: %s", ErrUsernameNotExists, username)
}

func FailedValidationError(errs error) error {
	return fmt.Errorf("%w: %w", ErrFailedValidation, errs)
}
