package hcs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// EncodeU4 generates the HCS-U4 code from normalized profile and CHIP
// This is a stub implementation using base64 encoding
// Can be replaced with base62 or other encoding schemes later
func EncodeU4(normalized *NormalizedProfile, chip string) (string, error) {
	// Create a compact structure for U4 encoding
	u4Data := map[string]interface{}{
		"profile": normalized,
		"chip":    chip,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(u4Data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal U4 data: %w", err)
	}

	// Encode to base64 URL-safe encoding (no padding)
	encoded := base64.RawURLEncoding.EncodeToString(jsonData)

	// Prefix with HCS-U4 identifier
	return fmt.Sprintf("HCS-U4|%s", encoded), nil
}

// DecodeU4 decodes an HCS-U4 code back to its components
func DecodeU4(code string) (*NormalizedProfile, string, error) {
	// Check prefix
	if len(code) < 7 || code[:7] != "HCS-U4|" {
		return nil, "", fmt.Errorf("invalid HCS-U4 format")
	}

	// Extract encoded part
	encoded := code[7:]

	// Decode from base64
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode U4: %w", err)
	}

	// Unmarshal JSON
	var u4Data struct {
		Profile *NormalizedProfile `json:"profile"`
		Chip    string             `json:"chip"`
	}

	if err := json.Unmarshal(decoded, &u4Data); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal U4 data: %w", err)
	}

	return u4Data.Profile, u4Data.Chip, nil
}
