package tests

import (
	"encoding/hex"
	"testing"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

func TestNormalizeProfile(t *testing.T) {
	input := &hcs.InputProfile{
		DominantElement: "Air",
		Modal: hcs.ModalBalance{
			Cardinal: 0.31,
			Fixed:    0.23,
			Mutable:  0.46,
		},
		Cognition: hcs.CognitionProfile{
			Fluid:        0.52,
			Crystallized: 0.13,
			Verbal:       0.53,
			Strategic:    0.15,
			Creative:     0.33,
		},
		Interaction: hcs.InteractionPreferences{
			Pace:      "balanced",
			Structure: "medium",
			Tone:      "precise",
		},
	}

	normalized := hcs.NormalizeProfile(input)

	// Check element mapping
	if normalized.Element != "A" {
		t.Errorf("Element should be 'A', got %s", normalized.Element)
	}

	// Check modal normalization
	if normalized.Modal.C != 31 {
		t.Errorf("Modal.Cardinal should be 31, got %d", normalized.Modal.C)
	}
	if normalized.Modal.F != 23 {
		t.Errorf("Modal.Fixed should be 23, got %d", normalized.Modal.F)
	}
	if normalized.Modal.M != 46 {
		t.Errorf("Modal.Mutable should be 46, got %d", normalized.Modal.M)
	}

	// Check cognition normalization
	if normalized.Cog.F != 52 {
		t.Errorf("Cog.Fluid should be 52, got %d", normalized.Cog.F)
	}
	if normalized.Cog.C != 13 {
		t.Errorf("Cog.Crystallized should be 13, got %d", normalized.Cog.C)
	}
	if normalized.Cog.V != 53 {
		t.Errorf("Cog.Verbal should be 53, got %d", normalized.Cog.V)
	}
	if normalized.Cog.S != 15 {
		t.Errorf("Cog.Strategic should be 15, got %d", normalized.Cog.S)
	}
	if normalized.Cog.Cr != 33 {
		t.Errorf("Cog.Creative should be 33, got %d", normalized.Cog.Cr)
	}

	// Check interaction mapping
	if normalized.Int.PB != "B" {
		t.Errorf("Int.PB should be 'B', got %s", normalized.Int.PB)
	}
	if normalized.Int.SM != "M" {
		t.Errorf("Int.SM should be 'M', got %s", normalized.Int.SM)
	}
	if normalized.Int.TN != "P" {
		t.Errorf("Int.TN should be 'P', got %s", normalized.Int.TN)
	}
}

func TestGenerateCHIP(t *testing.T) {
	// Create test salt
	salt := []byte("test-salt-32-bytes-padded-here!!")

	normalized := &hcs.NormalizedProfile{
		Element: "A",
		Modal: hcs.NormalizedModal{
			C: 31,
			F: 23,
			M: 46,
		},
		Cog: hcs.NormalizedCognition{
			F:  52,
			C:  13,
			V:  53,
			S:  15,
			Cr: 33,
		},
		Int: hcs.NormalizedInteraction{
			PB: "B",
			SM: "M",
			TN: "P",
		},
	}

	chip, err := hcs.GenerateCHIP(salt, normalized)
	if err != nil {
		t.Fatalf("GenerateCHIP failed: %v", err)
	}

	// Verify CHIP is 12 hex characters
	if len(chip) != 12 {
		t.Errorf("CHIP should be 12 characters, got %d", len(chip))
	}

	// Verify it's valid hex
	if _, err := hex.DecodeString(chip); err != nil {
		t.Errorf("CHIP is not valid hex: %s", chip)
	}

	t.Logf("Generated CHIP: %s", chip)

	// Generate again with same input - should be deterministic
	chip2, err := hcs.GenerateCHIP(salt, normalized)
	if err != nil {
		t.Fatalf("Second GenerateCHIP failed: %v", err)
	}

	if chip != chip2 {
		t.Errorf("CHIP generation not deterministic: %s != %s", chip, chip2)
	}

	// Different salt should produce different CHIP
	salt2 := []byte("different-salt-32-bytes-padded!!")
	chip3, err := hcs.GenerateCHIP(salt2, normalized)
	if err != nil {
		t.Fatalf("GenerateCHIP with different salt failed: %v", err)
	}

	if chip == chip3 {
		t.Errorf("Different salt should produce different CHIP: %s == %s", chip, chip3)
	}

	t.Logf("CHIP with different salt: %s", chip3)
}

func TestClampingBehavior(t *testing.T) {
	tests := []struct {
		name  string
		input *hcs.InputProfile
	}{
		{
			name: "values above 1.0",
			input: &hcs.InputProfile{
				DominantElement: "Fire",
				Modal: hcs.ModalBalance{
					Cardinal: 1.5,
					Fixed:    2.0,
					Mutable:  10.0,
				},
				Cognition: hcs.CognitionProfile{
					Fluid:        1.1,
					Crystallized: 1.2,
					Verbal:       1.3,
					Strategic:    1.4,
					Creative:     1.5,
				},
				Interaction: hcs.InteractionPreferences{
					Pace:      "fast",
					Structure: "high",
					Tone:      "sharp",
				},
			},
		},
		{
			name: "negative values",
			input: &hcs.InputProfile{
				DominantElement: "Water",
				Modal: hcs.ModalBalance{
					Cardinal: -0.5,
					Fixed:    -1.0,
					Mutable:  -10.0,
				},
				Cognition: hcs.CognitionProfile{
					Fluid:        -0.1,
					Crystallized: -0.2,
					Verbal:       -0.3,
					Strategic:    -0.4,
					Creative:     -0.5,
				},
				Interaction: hcs.InteractionPreferences{
					Pace:      "slow",
					Structure: "low",
					Tone:      "warm",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := hcs.NormalizeProfile(tt.input)

			// Check all modal values are clamped between 0 and 100
			if normalized.Modal.C < 0 || normalized.Modal.C > 100 {
				t.Errorf("Modal.C not clamped properly: %d", normalized.Modal.C)
			}
			if normalized.Modal.F < 0 || normalized.Modal.F > 100 {
				t.Errorf("Modal.F not clamped properly: %d", normalized.Modal.F)
			}
			if normalized.Modal.M < 0 || normalized.Modal.M > 100 {
				t.Errorf("Modal.M not clamped properly: %d", normalized.Modal.M)
			}

			// Check all cognition values are clamped
			if normalized.Cog.F < 0 || normalized.Cog.F > 100 {
				t.Errorf("Cog.F not clamped properly: %d", normalized.Cog.F)
			}
			if normalized.Cog.C < 0 || normalized.Cog.C > 100 {
				t.Errorf("Cog.C not clamped properly: %d", normalized.Cog.C)
			}
			if normalized.Cog.V < 0 || normalized.Cog.V > 100 {
				t.Errorf("Cog.V not clamped properly: %d", normalized.Cog.V)
			}
			if normalized.Cog.S < 0 || normalized.Cog.S > 100 {
				t.Errorf("Cog.S not clamped properly: %d", normalized.Cog.S)
			}
			if normalized.Cog.Cr < 0 || normalized.Cog.Cr > 100 {
				t.Errorf("Cog.Cr not clamped properly: %d", normalized.Cog.Cr)
			}

			t.Logf("%s - Clamped values: Modal(%d,%d,%d) Cog(%d,%d,%d,%d,%d)",
				tt.name,
				normalized.Modal.C, normalized.Modal.F, normalized.Modal.M,
				normalized.Cog.F, normalized.Cog.C, normalized.Cog.V, normalized.Cog.S, normalized.Cog.Cr)
		})
	}
}

func TestInteractionMappings(t *testing.T) {
	tests := []struct {
		pace      string
		structure string
		tone      string
		wantPB    string
		wantSM    string
		wantTN    string
	}{
		{"balanced", "low", "warm", "B", "L", "W"},
		{"fast", "medium", "neutral", "F", "M", "N"},
		{"slow", "high", "sharp", "S", "H", "S"},
		{"balanced", "medium", "precise", "B", "M", "P"},
		{"unknown", "invalid", "wrong", "B", "M", "N"}, // Test defaults
	}

	for _, tt := range tests {
		input := &hcs.InputProfile{
			DominantElement: "Earth",
			Modal:           hcs.ModalBalance{Cardinal: 0.5, Fixed: 0.3, Mutable: 0.2},
			Cognition:       hcs.CognitionProfile{Fluid: 0.5, Crystallized: 0.5, Verbal: 0.5, Strategic: 0.5, Creative: 0.5},
			Interaction: hcs.InteractionPreferences{
				Pace:      tt.pace,
				Structure: tt.structure,
				Tone:      tt.tone,
			},
		}

		normalized := hcs.NormalizeProfile(input)

		if normalized.Int.PB != tt.wantPB {
			t.Errorf("Pace %s should map to %s, got %s", tt.pace, tt.wantPB, normalized.Int.PB)
		}
		if normalized.Int.SM != tt.wantSM {
			t.Errorf("Structure %s should map to %s, got %s", tt.structure, tt.wantSM, normalized.Int.SM)
		}
		if normalized.Int.TN != tt.wantTN {
			t.Errorf("Tone %s should map to %s, got %s", tt.tone, tt.wantTN, normalized.Int.TN)
		}
	}
}
