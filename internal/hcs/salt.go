package hcs

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

const (
	saltFileName = ".hcs_salt"
	saltSize     = 32
)

// LoadOrCreateSalt loads the salt from file or creates a new one if not exists
func LoadOrCreateSalt(dir string) ([]byte, error) {
	if dir == "" {
		dir = "."
	}

	saltPath := filepath.Join(dir, saltFileName)

	// Try to read existing salt
	salt, err := os.ReadFile(saltPath)
	if err == nil && len(salt) == saltSize {
		return salt, nil
	}

	// Generate new salt if file doesn't exist or is invalid
	if os.IsNotExist(err) || len(salt) != saltSize {
		salt = make([]byte, saltSize)
		if _, err := rand.Read(salt); err != nil {
			return nil, fmt.Errorf("failed to generate salt: %w", err)
		}

		// Write salt to file
		if err := os.WriteFile(saltPath, salt, 0600); err != nil {
			return nil, fmt.Errorf("failed to save salt: %w", err)
		}

		return salt, nil
	}

	return nil, fmt.Errorf("failed to read salt: %w", err)
}
