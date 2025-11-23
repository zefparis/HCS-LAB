package hcs

import (
	"encoding/hex"
	"fmt"
	"os"
	"sync"
)

var (
	secretKeyOnce sync.Once
	secretKey     []byte
	secretKeyErr  error
)

// LoadSecretKey loads and caches the HCS secret key from the HCS_SECRET_KEY environment variable.
// The key must be a hex string representing either 32 or 64 bytes. The raw key material is never logged.
func LoadSecretKey() ([]byte, error) {
	secretKeyOnce.Do(func() {
		value := os.Getenv("HCS_SECRET_KEY")
		if value == "" {
			secretKeyErr = fmt.Errorf("HCS_SECRET_KEY is not set")
			return
		}

		decoded, err := hex.DecodeString(value)
		if err != nil {
			secretKeyErr = fmt.Errorf("invalid HCS_SECRET_KEY hex encoding: %w", err)
			return
		}

		if l := len(decoded); l != 32 && l != 64 {
			secretKeyErr = fmt.Errorf("HCS_SECRET_KEY must be 32 or 64 bytes, got %d bytes", l)
			return
		}

		secretKey = decoded
	})

	return secretKey, secretKeyErr
}
