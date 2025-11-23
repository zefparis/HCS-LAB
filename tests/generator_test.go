package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

func getTestInput() *hcs.InputProfile {
	return &hcs.InputProfile{
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
}

func TestGeneratorDeterminism(t *testing.T) {
	// Ensure secret key is set so that U7 generation can complete
	if err := os.Setenv("HCS_SECRET_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"); err != nil {
		t.Fatalf("failed to set HCS_SECRET_KEY: %v", err)
	}
	// Create a temporary directory for salt
	tempDir := t.TempDir()

	// Create generator with specific salt directory
	gen1, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Generate first output
	input := getTestInput()
	output1, err := gen1.Generate(input)
	if err != nil {
		t.Fatalf("Failed to generate output1: %v", err)
	}

	// Create another generator instance with same salt directory
	gen2, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to create generator2: %v", err)
	}

	// Generate second output with same input
	output2, err := gen2.Generate(input)
	if err != nil {
		t.Fatalf("Failed to generate output2: %v", err)
	}

	// Check that outputs are identical
	if output1.CodeU3 != output2.CodeU3 {
		t.Errorf("CodeU3 not deterministic: %s != %s", output1.CodeU3, output2.CodeU3)
	}
	if output1.CodeU4 != output2.CodeU4 {
		t.Errorf("CodeU4 not deterministic: %s != %s", output1.CodeU4, output2.CodeU4)
	}
	if output1.Chip != output2.Chip {
		t.Errorf("Chip not deterministic: %s != %s", output1.Chip, output2.Chip)
	}

	t.Logf("Generated HCS-U3: %s", output1.CodeU3)
	t.Logf("Generated CHIP: %s", output1.Chip)
}

func TestSaltPersistence(t *testing.T) {
	if err := os.Setenv("HCS_SECRET_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"); err != nil {
		t.Fatalf("failed to set HCS_SECRET_KEY: %v", err)
	}
	// Create a temporary directory
	tempDir := t.TempDir()
	saltPath := filepath.Join(tempDir, ".hcs_salt")

	// Generate with first salt
	gen1, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to create generator1: %v", err)
	}

	input := getTestInput()
	output1, err := gen1.Generate(input)
	if err != nil {
		t.Fatalf("Failed to generate output1: %v", err)
	}

	// Delete salt file
	if err := os.Remove(saltPath); err != nil {
		t.Fatalf("Failed to remove salt file: %v", err)
	}

	// Generate with new salt
	gen2, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to create generator2: %v", err)
	}

	output2, err := gen2.Generate(input)
	if err != nil {
		t.Fatalf("Failed to generate output2: %v", err)
	}

	// Chips should be different after salt regeneration
	if output1.Chip == output2.Chip {
		t.Errorf("Chip should be different after salt regeneration: %s == %s", output1.Chip, output2.Chip)
	}

	t.Logf("Original CHIP: %s", output1.Chip)
	t.Logf("New CHIP after salt regeneration: %s", output2.Chip)
}

func TestGeneratorValidation(t *testing.T) {
	tempDir := t.TempDir()
	gen, err := hcs.NewGeneratorWithSaltDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	tests := []struct {
		name    string
		input   *hcs.InputProfile
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   getTestInput(),
			wantErr: false,
		},
		{
			name: "invalid element",
			input: &hcs.InputProfile{
				DominantElement: "Invalid",
				Modal:           hcs.ModalBalance{Cardinal: 0.5, Fixed: 0.3, Mutable: 0.2},
				Cognition:       hcs.CognitionProfile{Fluid: 0.5, Crystallized: 0.5, Verbal: 0.5, Strategic: 0.5, Creative: 0.5},
				Interaction:     hcs.InteractionPreferences{Pace: "balanced", Structure: "medium", Tone: "neutral"},
			},
			wantErr: true,
		},
		{
			name: "out of range modal",
			input: &hcs.InputProfile{
				DominantElement: "Earth",
				Modal:           hcs.ModalBalance{Cardinal: 1.5, Fixed: 0.3, Mutable: 0.2},
				Cognition:       hcs.CognitionProfile{Fluid: 0.5, Crystallized: 0.5, Verbal: 0.5, Strategic: 0.5, Creative: 0.5},
				Interaction:     hcs.InteractionPreferences{Pace: "balanced", Structure: "medium", Tone: "neutral"},
			},
			wantErr: true,
		},
		{
			name: "invalid pace",
			input: &hcs.InputProfile{
				DominantElement: "Water",
				Modal:           hcs.ModalBalance{Cardinal: 0.5, Fixed: 0.3, Mutable: 0.2},
				Cognition:       hcs.CognitionProfile{Fluid: 0.5, Crystallized: 0.5, Verbal: 0.5, Strategic: 0.5, Creative: 0.5},
				Interaction:     hcs.InteractionPreferences{Pace: "invalid", Structure: "medium", Tone: "neutral"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gen.Generate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
