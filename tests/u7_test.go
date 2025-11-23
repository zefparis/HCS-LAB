package tests

import (
	"os"
	"testing"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

func setTestSecretKey(t *testing.T) {
	t.Helper()
	// 32-byte key (64 hex chars)
	if err := os.Setenv("HCS_SECRET_KEY", "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"); err != nil {
		t.Fatalf("failed to set HCS_SECRET_KEY: %v", err)
	}
}

// TestU7Determinism verifies that with a fixed secret, salt, and profile,
// U7-related outputs are fully deterministic.
func TestU7Determinism(t *testing.T) {
	setTestSecretKey(t)

	tempDir := t.TempDir()
	gen1, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("failed to create generator1: %v", err)
	}

	input := getTestInput()
	out1, err := gen1.Generate(input)
	if err != nil {
		t.Fatalf("failed to generate out1: %v", err)
	}

	gen2, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("failed to create generator2: %v", err)
	}

	out2, err := gen2.Generate(input)
	if err != nil {
		t.Fatalf("failed to generate out2: %v", err)
	}

	if out1.CodeU7 != out2.CodeU7 {
		t.Errorf("CodeU7 not deterministic: %s != %s", out1.CodeU7, out2.CodeU7)
	}
	if out1.QSig != out2.QSig {
		t.Errorf("QSIG not deterministic: %s != %s", out1.QSig, out2.QSig)
	}
	if out1.B3Sig != out2.B3Sig {
		t.Errorf("B3Sig not deterministic: %s != %s", out1.B3Sig, out2.B3Sig)
	}
}

// TestU7SecretSensitivity ensures different secrets produce different signatures
// and thus different U7 codes for the same input and salt.
func TestU7SecretSensitivity(t *testing.T) {
	// First secret
	if err := os.Setenv("HCS_SECRET_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"); err != nil {
		t.Fatalf("failed to set secret1: %v", err)
	}

	tempDir := t.TempDir()
	gen1, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("failed to create generator1: %v", err)
	}

	input := getTestInput()
	out1, err := gen1.Generate(input)
	if err != nil {
		t.Fatalf("failed to generate out1: %v", err)
	}

	// Second secret (same salt dir)
	if err := os.Setenv("HCS_SECRET_KEY", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"); err != nil {
		t.Fatalf("failed to set secret2: %v", err)
	}

	gen2, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("failed to create generator2: %v", err)
	}

	out2, err := gen2.Generate(input)
	if err != nil {
		t.Fatalf("failed to generate out2: %v", err)
	}

	if out1.QSig == out2.QSig {
		t.Errorf("QSIG should differ for different secrets")
	}
	if out1.B3Sig == out2.B3Sig {
		t.Errorf("B3Sig should differ for different secrets")
	}
	if out1.CodeU7 == out2.CodeU7 {
		t.Errorf("CodeU7 should differ for different secrets")
	}
}

// TestU7Avalanche performs a basic avalanche sanity check: a small change
// in the input should cause a completely different signature.
func TestU7Avalanche(t *testing.T) {
	setTestSecretKey(t)

	tempDir := t.TempDir()
	gen, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	base := getTestInput()
	out1, err := gen.Generate(base)
	if err != nil {
		t.Fatalf("failed to generate base output: %v", err)
	}

	// Slightly tweak cognition
	mutated := getTestInput()
	mutated.Cognition.Fluid += 0.01

	out2, err := gen.Generate(mutated)
	if err != nil {
		t.Fatalf("failed to generate mutated output: %v", err)
	}

	if out1.QSig == out2.QSig {
		t.Errorf("QSIG should change when profile is slightly modified")
	}
}
