package utils

import (
	"crypto/ed25519"
	"encoding/base64"
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

func GetTokenFromRequest(req *http.Request) string {
	return req.Header.Get("Authorization")
}

func DecodeJWTSecret(jwtSecret string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(jwtSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("failed decoding jwt secret from envs: %w", err)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, nil, fmt.Errorf("invalid private key length, expected (%d) got (%d)", ed25519.PrivateKeySize, len(privateKeyBytes))
	}

	newPriv := ed25519.NewKeyFromSeed(privateKeyBytes[:32])

	for i, b := range privateKeyBytes {
		if b != newPriv[i] {
			return nil, nil, fmt.Errorf("provided public key does not match derived public key")
		}
	}

	priv := ed25519.PrivateKey(privateKeyBytes)
	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("failed asserting crypto public key to ed25519 public key")
	}

	return priv, pub, nil
}
