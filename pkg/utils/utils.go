package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func ReadRequestJSON(req *http.Request, v any) error {
	if req.Body == nil {
		return fmt.Errorf("missing request body")
	}

	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return fmt.Errorf("failed decoding request body: %w", err)
	}

	return nil
}

func WriteJSON(rw http.ResponseWriter, status int, v any) error {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)

	if err := json.NewEncoder(rw).Encode(v); err != nil {
		return fmt.Errorf("failed encoding response json: %w", err)
	}

	return nil
}

func WriteError(rw http.ResponseWriter, status int) {
	http.Error(rw, http.StatusText(status), status)
}

func WriteCustomError(rw http.ResponseWriter, status int, err error) {
	http.Error(rw, fmt.Sprintf("%s: %s", http.StatusText(status), err.Error()), status)
}
