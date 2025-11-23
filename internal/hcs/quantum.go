package hcs

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"

	"github.com/zeebo/blake3"
	"golang.org/x/crypto/sha3"
)

// ComputeQuantumSignatures computes the primary HMAC-SHA3-256 based signature (QSIG)
// and a secondary BLAKE3 digest over the canonical profile data. The design is
// intentionally conservative: 256-bit symmetric hashes remain extremely strong
// even under generic quantum attacks such as Grover's algorithm.
func ComputeQuantumSignatures(canonical []byte, secret []byte, salt []byte) (qsigHex string, b3Hex string, err error) {
	if len(secret) == 0 {
		return "", "", fmt.Errorf("secret key must not be empty")
	}

	// Derive a per-instance key from the master secret and salt using HMAC-SHA3-256.
	// salt is treated as public diversification material; the secret remains private.
	h := hmac.New(sha3.New256, secret)
	if _, err = h.Write(salt); err != nil {
		return "", "", fmt.Errorf("failed to derive key: %w", err)
	}
	derivedKey := h.Sum(nil)

	// Primary quantum-style signature: HMAC-SHA3-256 over canonical profile.
	h2 := hmac.New(sha3.New256, derivedKey)
	if _, err = h2.Write(canonical); err != nil {
		return "", "", fmt.Errorf("failed to compute primary signature: %w", err)
	}
	qsig := h2.Sum(nil)
	qsigHex = hex.EncodeToString(qsig)

	// Secondary digest: BLAKE3 over canonical profile (optionally salted via derived key).
	b3 := blake3.New()
	if _, err = b3.Write(derivedKey); err != nil {
		return "", "", fmt.Errorf("failed to update BLAKE3 with key: %w", err)
	}
	if _, err = b3.Write(canonical); err != nil {
		return "", "", fmt.Errorf("failed to update BLAKE3 with canonical: %w", err)
	}
	b3sum := b3.Sum(nil)
	b3Hex = hex.EncodeToString(b3sum)

	return qsigHex, b3Hex, nil
}
