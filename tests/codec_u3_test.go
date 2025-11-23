package tests

import (
	"regexp"
	"testing"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

func TestEncodeU3Format(t *testing.T) {
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

	chip := "b04edb83f10e"
	result := hcs.EncodeU3(input, chip)

	// Expected format
	expected := "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb83f10e"

	if result != expected {
		t.Errorf("EncodeU3() = %s, want %s", result, expected)
	}

	// Validate format using regex
	if !hcs.ValidateU3Format(result) {
		t.Errorf("Generated U3 code failed format validation: %s", result)
	}

	t.Logf("Generated U3: %s", result)
}

func TestValidateU3Format(t *testing.T) {
	tests := []struct {
		name  string
		code  string
		valid bool
	}{
		{
			name:  "valid U3 code",
			code:  "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb83f10e",
			valid: true,
		},
		{
			name:  "valid with Earth element",
			code:  "HCS-U3|E:E|MOD:c00f99m01|COG:F00C99V50S50Cr50|INT:PB=F,SM=L,TN=W|CHIP:abc123def456",
			valid: true,
		},
		{
			name:  "valid with Water element",
			code:  "HCS-U3|E:W|MOD:c50f25m25|COG:F33C33V33S00Cr00|INT:PB=S,SM=H,TN=S|CHIP:123456789abc",
			valid: true,
		},
		{
			name:  "valid with Fire element",
			code:  "HCS-U3|E:F|MOD:c33f33m34|COG:F25C25V25S25Cr00|INT:PB=B,SM=M,TN=N|CHIP:fedcba987654",
			valid: true,
		},
		{
			name:  "invalid - wrong prefix",
			code:  "HCS-U2|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb83f10e",
			valid: false,
		},
		{
			name:  "invalid - invalid element",
			code:  "HCS-U3|E:X|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb83f10e",
			valid: false,
		},
		{
			name:  "invalid - missing modal values",
			code:  "HCS-U3|E:A|MOD:c31f23|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb83f10e",
			valid: false,
		},
		{
			name:  "invalid - wrong chip length",
			code:  "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb",
			valid: false,
		},
		{
			name:  "invalid - invalid pace",
			code:  "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=X,SM=M,TN=P|CHIP:b04edb83f10e",
			valid: false,
		},
		{
			name:  "invalid - invalid structure",
			code:  "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=X,TN=P|CHIP:b04edb83f10e",
			valid: false,
		},
		{
			name:  "invalid - invalid tone",
			code:  "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=X|CHIP:b04edb83f10e",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hcs.ValidateU3Format(tt.code)
			if result != tt.valid {
				t.Errorf("ValidateU3Format(%s) = %v, want %v", tt.code, result, tt.valid)
			}
		})
	}
}

func TestU3Regex(t *testing.T) {
	// Test the regex pattern directly
	pattern := `^HCS-U3\|E:[AEWF]\|MOD:c\d{2}f\d{2}m\d{2}\|COG:F\d{2}C\d{2}V\d{2}S\d{2}Cr\d{2}\|INT:PB=[BFS],SM=[LMH],TN=[WNSP]\|CHIP:[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)

	validCode := "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:b04edb83f10e"

	if !re.MatchString(validCode) {
		t.Errorf("Regex should match valid code: %s", validCode)
	}

	t.Logf("Regex pattern: %s", pattern)
	t.Logf("Test code: %s", validCode)
}

func TestElementMapping(t *testing.T) {
	tests := []struct {
		element string
		want    string
	}{
		{"Earth", "E"},
		{"Air", "A"},
		{"Water", "W"},
		{"Fire", "F"},
	}

	for _, tt := range tests {
		input := &hcs.InputProfile{
			DominantElement: tt.element,
			Modal:           hcs.ModalBalance{Cardinal: 0.5, Fixed: 0.3, Mutable: 0.2},
			Cognition:       hcs.CognitionProfile{Fluid: 0.5, Crystallized: 0.5, Verbal: 0.5, Strategic: 0.5, Creative: 0.5},
			Interaction:     hcs.InteractionPreferences{Pace: "balanced", Structure: "medium", Tone: "neutral"},
		}

		result := hcs.EncodeU3(input, "123456789abc")

		// Check if element is correctly encoded
		expectedPrefix := "HCS-U3|E:" + tt.want
		if len(result) < len(expectedPrefix) || result[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("Element %s should encode to %s, got %s", tt.element, tt.want, result)
		}
	}
}

func TestPercentageRounding(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"zero", 0.0, "00"},
		{"small", 0.004, "00"},
		{"round down", 0.314, "31"},
		{"exact", 0.31, "31"},
		{"round up", 0.315, "32"},
		{"half", 0.5, "50"},
		{"large", 0.994, "99"},
		{"one", 1.0, "100"},       // Note: This becomes "100" which is 3 digits
		{"above one", 1.5, "100"}, // Clamped to 1.0
		{"negative", -0.5, "00"},  // Clamped to 0
	}

	for _, tt := range tests {
		input := &hcs.InputProfile{
			DominantElement: "Earth",
			Modal: hcs.ModalBalance{
				Cardinal: tt.value,
				Fixed:    0.5,
				Mutable:  0.5,
			},
			Cognition:   hcs.CognitionProfile{Fluid: 0.5, Crystallized: 0.5, Verbal: 0.5, Strategic: 0.5, Creative: 0.5},
			Interaction: hcs.InteractionPreferences{Pace: "balanced", Structure: "medium", Tone: "neutral"},
		}

		result := hcs.EncodeU3(input, "123456789abc")

		// Extract the cardinal value from the result
		// Format: HCS-U3|E:E|MOD:cXXf50m50|...
		re := regexp.MustCompile(`MOD:c(\d+)f`)
		matches := re.FindStringSubmatch(result)
		if len(matches) < 2 {
			t.Errorf("Failed to extract cardinal value from: %s", result)
			continue
		}

		// Handle the special case where 1.0 becomes "100"
		expectedResult := tt.expected
		if tt.value >= 1.0 {
			expectedResult = "100"
		}

		// We expect a 2-digit format in the result for values < 100
		if len(expectedResult) == 3 {
			expectedResult = "100"
		} else if len(expectedResult) < 2 {
			expectedResult = "0" + expectedResult
		}

		// The actual result will be either 2 or 3 digits
		if matches[1] != expectedResult && !(expectedResult == "100" && matches[1] == "100") {
			t.Errorf("%s: value %f should round to %s, got %s", tt.name, tt.value, expectedResult[:2], matches[1])
		}
	}
}
