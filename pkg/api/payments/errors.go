package payments

import (
	"errors"
	"fmt"
)

var (
	ErrAccountNotExists          = errors.New("account not found for user and currency")
	ErrRecipientAccountNotExists = errors.New("the recipient does not have an account in that currency")
	ErrNotEnoughFunds            = errors.New("not enough funds to process the transaction")
	ErrRecipientNotExsits        = errors.New("the recipient does not exist")
)

func ErrFailedValidation(errs error) error {
	return fmt.Errorf("failed validation: %w", errs)
}
