package payments

import (
	"fmt"
)

func ErrFailedValidation(errs error) error {
	return fmt.Errorf("failed validation: %w", errs)
}
