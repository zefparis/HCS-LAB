package tests

import (
	"strings"
	"testing"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

// TestEncodeU5 tests the HCS-U5 encoding
func TestEncodeU5(t *testing.T) {
	// Create sample profiles
	western := &hcs.WesternProfile{
		DominantElement: "Fire",
		Modal: hcs.ModalBalance{
			Cardinal: 0.4,
			Fixed:    0.3,
			Mutable:  0.3,
		},
		Cognition: hcs.CognitionProfile{
			Fluid:        0.7,
			Crystallized: 0.6,
			Verbal:       0.5,
			Strategic:    0.8,
			Creative:     0.9,
		},
		Interaction: hcs.InteractionPreferences{
			Pace:      "fast",
			Structure: "medium",
			Tone:      "sharp",
		},
	}

	chinese := &hcs.ChineseProfile{
		YearPillar:     "Geng-Wu",
		MonthPillar:    "Ding-Si",
		DayPillar:      "Jia-Chen",
		HourPillar:     "Xin-Wei",
		YinYangBalance: 0.6,
		ElementBalance: map[string]float64{
			"Wood":  0.2,
			"Fire":  0.3,
			"Earth": 0.2,
			"Metal": 0.15,
			"Water": 0.15,
		},
		DayMaster:         "Jia",
		DayMasterStrength: 0.7,
	}

	fusion := hcs.BuildFusionProfile(western, chinese)

	// Generate salt for testing
	salt := []byte("test-salt-12345")

	// Encode U5
	u5Code, err := hcs.EncodeU5(western, chinese, fusion, salt)
	if err != nil {
		t.Fatalf("Failed to encode U5: %v", err)
	}

	// Check format
	if !strings.HasPrefix(u5Code, "HCS-U5|") {
		t.Error("U5 code doesn't start with HCS-U5|")
	}

	// Validate format
	if !hcs.ValidateU5Format(u5Code) {
		t.Errorf("Generated U5 code fails validation: %s", u5Code)
	}

	// Check required segments
	requiredSegments := []string{"|W:", "|C:", "|F:", "|CHIP:"}
	for _, seg := range requiredSegments {
		if !strings.Contains(u5Code, seg) {
			t.Errorf("U5 code missing segment %s", seg)
		}
	}

	// Test determinism - same input should produce same output
	u5Code2, err := hcs.EncodeU5(western, chinese, fusion, salt)
	if err != nil {
		t.Fatalf("Failed to encode U5 second time: %v", err)
	}

	if u5Code != u5Code2 {
		t.Error("U5 encoding is not deterministic")
	}

	// Different salt should produce different CHIP
	salt2 := []byte("different-salt-67890")
	u5Code3, err := hcs.EncodeU5(western, chinese, fusion, salt2)
	if err != nil {
		t.Fatalf("Failed to encode U5 with different salt: %v", err)
	}

	// Extract CHIPs
	chip1 := extractChip(u5Code)
	chip3 := extractChip(u5Code3)

	if chip1 == chip3 {
		t.Error("Different salts produced same CHIP")
	}
}

// TestValidateU5Format tests the U5 format validation
func TestValidateU5Format(t *testing.T) {
	testCases := []struct {
		code     string
		expected bool
		desc     string
	}{
		{
			"HCS-U5|A1|W:1234|C:5678|F:9abc|CHIP:def012345678",
			true,
			"Valid U5 format",
		},
		{
			"HCS-U3|A1|W:1234|C:5678|F:9abc|CHIP:def012345678",
			false,
			"Wrong prefix",
		},
		{
			"HCS-U5|A1|W:1234|C:5678|CHIP:def012345678",
			false,
			"Missing F segment",
		},
		{
			"HCS-U5|A1|W:1234|C:5678|F:9abc|CHIP:def",
			true, // Chip can be shorter in validation
			"Short CHIP",
		},
		{
			"",
			false,
			"Empty string",
		},
		{
			"HCS-U5",
			false,
			"Incomplete code",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := hcs.ValidateU5Format(tc.code)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.desc, result)
			}
		})
	}
}

// TestDecodeU5 tests the U5 decoding
func TestDecodeU5(t *testing.T) {
	validCode := "HCS-U5|A1|W:1234|C:5678|F:9abc|CHIP:def012345678"

	components, err := hcs.DecodeU5(validCode)
	if err != nil {
		t.Fatalf("Failed to decode valid U5: %v", err)
	}

	// Check fusion ID
	if components["fusionId"] != "A1" {
		t.Errorf("Expected fusion ID 'A1', got '%s'", components["fusionId"])
	}

	// Check western hex
	if components["western"] != "1234" {
		t.Errorf("Expected western hex '1234', got '%s'", components["western"])
	}

	// Check chinese hex
	if components["chinese"] != "5678" {
		t.Errorf("Expected chinese hex '5678', got '%s'", components["chinese"])
	}

	// Check fusion hex
	if components["fusion"] != "9abc" {
		t.Errorf("Expected fusion hex '9abc', got '%s'", components["fusion"])
	}

	// Check CHIP
	if components["chip"] != "def012345678" {
		t.Errorf("Expected CHIP 'def012345678', got '%s'", components["chip"])
	}

	// Test invalid code
	invalidCode := "HCS-U3|invalid"
	_, err = hcs.DecodeU5(invalidCode)
	if err == nil {
		t.Error("Decoding invalid U5 should return error")
	}
}

// TestFusionProfile tests the fusion profile generation
func TestFusionProfile(t *testing.T) {
	western := &hcs.WesternProfile{
		DominantElement: "Fire",
		Modal: hcs.ModalBalance{
			Cardinal: 0.5,
			Fixed:    0.3,
			Mutable:  0.2,
		},
		Cognition: hcs.CognitionProfile{
			Fluid:        0.8,
			Crystallized: 0.6,
			Verbal:       0.7,
			Strategic:    0.9,
			Creative:     0.85,
		},
		Interaction: hcs.InteractionPreferences{
			Pace:      "fast",
			Structure: "high",
			Tone:      "precise",
		},
	}

	chinese := &hcs.ChineseProfile{
		YearPillar:     "Bing-Yin",
		MonthPillar:    "Wu-Chen",
		DayPillar:      "Ren-Wu",
		HourPillar:     "Gui-Hai",
		YinYangBalance: 0.7,
		ElementBalance: map[string]float64{
			"Wood":  0.1,
			"Fire":  0.4,
			"Earth": 0.2,
			"Metal": 0.1,
			"Water": 0.2,
		},
		DayMaster:         "Ren",
		DayMasterStrength: 0.6,
	}

	fusion := hcs.BuildFusionProfile(western, chinese)

	// Check fusion profile has all required fields
	if fusion.FusionID == "" {
		t.Error("Fusion ID is empty")
	}

	if len(fusion.FusionID) != 2 {
		t.Errorf("Fusion ID should be 2 characters, got %d", len(fusion.FusionID))
	}

	// Check element signature
	if len(fusion.ElementSignature) == 0 {
		t.Error("Element signature is empty")
	}

	// Element signature should sum to 1
	total := 0.0
	for _, val := range fusion.ElementSignature {
		total += val
	}
	if total < 0.99 || total > 1.01 {
		t.Errorf("Element signature doesn't sum to 1.0: %f", total)
	}

	// Check cognitive fusion values are in range
	cogValues := []float64{
		fusion.CognitiveFusion.Analytical,
		fusion.CognitiveFusion.Creative,
		fusion.CognitiveFusion.Grounded,
		fusion.CognitiveFusion.Adaptive,
		fusion.CognitiveFusion.Expressive,
	}

	for i, val := range cogValues {
		if val < 0 || val > 1 {
			t.Errorf("Cognitive fusion value %d out of range: %f", i, val)
		}
	}

	// Check tempo signals
	if fusion.TempoSignals.Pace < 0 || fusion.TempoSignals.Pace > 1 {
		t.Errorf("Tempo pace out of range: %f", fusion.TempoSignals.Pace)
	}

	if fusion.TempoSignals.Rhythm == "" {
		t.Error("Tempo rhythm is empty")
	}

	// Check balance metrics
	if fusion.UnifiedBalance < 0 || fusion.UnifiedBalance > 1 {
		t.Errorf("Unified balance out of range: %f", fusion.UnifiedBalance)
	}

	if fusion.HarmonicResonance < 0 || fusion.HarmonicResonance > 1 {
		t.Errorf("Harmonic resonance out of range: %f", fusion.HarmonicResonance)
	}
}

// TestFullIntegration tests the complete flow with generator
func TestFullIntegration(t *testing.T) {
	// Create generator
	gen, err := hcs.NewGeneratorWithSaltDir(".")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Create input profile with birth info
	input := &hcs.InputProfile{
		DominantElement: "Water",
		Modal: hcs.ModalBalance{
			Cardinal: 0.3,
			Fixed:    0.5,
			Mutable:  0.2,
		},
		Cognition: hcs.CognitionProfile{
			Fluid:        0.75,
			Crystallized: 0.65,
			Verbal:       0.8,
			Strategic:    0.7,
			Creative:     0.6,
		},
		Interaction: hcs.InteractionPreferences{
			Pace:      "balanced",
			Structure: "medium",
			Tone:      "warm",
		},
		BirthInfo: &hcs.BirthInfo{
			Year:     1985,
			Month:    7,
			Day:      20,
			Hour:     10,
			Minute:   45,
			Timezone: "UTC",
		},
	}

	// Generate HCS codes
	output, err := gen.Generate(input)
	if err != nil {
		t.Fatalf("Failed to generate HCS: %v", err)
	}

	// Check U3 and U4 are still generated
	if output.CodeU3 == "" {
		t.Error("U3 code not generated")
	}

	if output.CodeU4 == "" {
		t.Error("U4 code not generated")
	}

	// Check new fields
	if output.CodeU5 == "" {
		t.Error("U5 code not generated")
	}

	if output.ChineseProfile == nil {
		t.Error("Chinese profile not generated")
	}

	if output.CombinedProfile == nil {
		t.Error("Combined profile not generated")
	}

	// Validate U5 format
	if !hcs.ValidateU5Format(output.CodeU5) {
		t.Errorf("Generated U5 code is invalid: %s", output.CodeU5)
	}

	// Test without birth info - should still work but without Chinese features
	inputNoBirth := &hcs.InputProfile{
		DominantElement: "Earth",
		Modal: hcs.ModalBalance{
			Cardinal: 0.4,
			Fixed:    0.4,
			Mutable:  0.2,
		},
		Cognition: hcs.CognitionProfile{
			Fluid:        0.5,
			Crystallized: 0.5,
			Verbal:       0.5,
			Strategic:    0.5,
			Creative:     0.5,
		},
		Interaction: hcs.InteractionPreferences{
			Pace:      "slow",
			Structure: "low",
			Tone:      "neutral",
		},
	}

	outputNoBirth, err := gen.Generate(inputNoBirth)
	if err != nil {
		t.Fatalf("Failed to generate HCS without birth info: %v", err)
	}

	// Should have U3 and U4 but not U5 or Chinese profile
	if outputNoBirth.CodeU3 == "" {
		t.Error("U3 not generated without birth info")
	}

	if outputNoBirth.CodeU5 != "" {
		t.Error("U5 should not be generated without birth info")
	}

	if outputNoBirth.ChineseProfile != nil {
		t.Error("Chinese profile should not be generated without birth info")
	}
}

// Helper function to extract CHIP from U5 code
func extractChip(code string) string {
	chipPrefix := "CHIP:"
	idx := strings.Index(code, chipPrefix)
	if idx == -1 {
		return ""
	}
	start := idx + len(chipPrefix)
	if start+12 > len(code) {
		return code[start:]
	}
	return code[start : start+12]
}
