package hcs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// EncodeU5 generates the HCS-U5 code from combined profiles and CHIP
// Format: HCS-U5|XX|W:<hex>|C:<hex>|F:<hex>|CHIP:<12hex>
func EncodeU5(western *WesternProfile, chinese *ChineseProfile, fusion *FusionProfile, salt []byte) (string, error) {
	// Generate fusion ID (2 chars)
	fusionID := fusion.FusionID
	if len(fusionID) != 2 {
		fusionID = "XX" // Fallback
	}

	// Compress Western profile to 16-bit hex
	westernHex := compressWesternProfile(western)

	// Compress Chinese profile to 16-bit hex
	chineseHex := compressChineseProfile(chinese)

	// Compress Fusion traits to 16-bit hex
	fusionHex := compressFusionProfile(fusion)

	// Generate CHIP for U5 (using combined data)
	chipU5, err := generateU5Chip(western, chinese, fusion, salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate U5 CHIP: %w", err)
	}

	// Construct HCS-U5 code
	code := fmt.Sprintf("HCS-U5|%s|W:%s|C:%s|F:%s|CHIP:%s",
		fusionID, westernHex, chineseHex, fusionHex, chipU5)

	return code, nil
}

// compressWesternProfile compresses Western profile to 4 hex chars (16 bits)
func compressWesternProfile(western *WesternProfile) string {
	// Allocate 16 bits:
	// - 2 bits: element (4 options)
	// - 3x3 bits: modal balance (each 0-7 scale)
	// - 2 bits: pace
	// - 2 bits: structure
	// - 1 bit: tone category

	var bits uint16

	// Element (bits 15-14)
	elementBits := uint16(0)
	switch western.DominantElement {
	case "Fire":
		elementBits = 0
	case "Earth":
		elementBits = 1
	case "Air":
		elementBits = 2
	case "Water":
		elementBits = 3
	}
	bits |= (elementBits << 14)

	// Modal balance (bits 13-5)
	// Convert to 3-bit values (0-7)
	cardinalBits := uint16(western.Modal.Cardinal * 7)
	fixedBits := uint16(western.Modal.Fixed * 7)
	mutableBits := uint16(western.Modal.Mutable * 7)

	bits |= (cardinalBits << 10)
	bits |= (fixedBits << 7)
	bits |= (mutableBits << 4)

	// Pace (bits 3-2)
	paceBits := uint16(0)
	switch western.Interaction.Pace {
	case "slow":
		paceBits = 0
	case "balanced":
		paceBits = 1
	case "fast":
		paceBits = 2
	}
	bits |= (paceBits << 2)

	// Structure (bit 1)
	if western.Interaction.Structure == "high" {
		bits |= (1 << 1)
	}

	// Tone category (bit 0)
	if western.Interaction.Tone == "sharp" || western.Interaction.Tone == "precise" {
		bits |= 1
	}

	// Convert to 4-char hex
	return fmt.Sprintf("%04x", bits)
}

// compressChineseProfile compresses Chinese profile to 4 hex chars (16 bits)
func compressChineseProfile(chinese *ChineseProfile) string {
	// Allocate 16 bits:
	// - 3 bits: dominant element (5 options + padding)
	// - 3 bits: yin/yang balance (0-7 scale)
	// - 4 bits: day master (10 stems)
	// - 3 bits: day master strength
	// - 3 bits: element distribution pattern

	var bits uint16

	// Dominant element (bits 15-13)
	dominantElement := chinese.GetDominantChineseElement()
	elementBits := uint16(0)
	switch dominantElement {
	case "Wood":
		elementBits = 0
	case "Fire":
		elementBits = 1
	case "Earth":
		elementBits = 2
	case "Metal":
		elementBits = 3
	case "Water":
		elementBits = 4
	}
	bits |= (elementBits << 13)

	// Yin/Yang balance (bits 12-10)
	yinYangBits := uint16(chinese.YinYangBalance * 7)
	bits |= (yinYangBits << 10)

	// Day Master index (bits 9-6)
	dayMasterBits := uint16(0)
	for i, stem := range HeavenlyStems {
		if stem.Name == chinese.DayMaster {
			dayMasterBits = uint16(i)
			break
		}
	}
	bits |= (dayMasterBits << 6)

	// Day Master strength (bits 5-3)
	strengthBits := uint16(chinese.DayMasterStrength * 7)
	bits |= (strengthBits << 3)

	// Element distribution pattern (bits 2-0)
	// Encode whether elements are balanced or skewed
	variance := calculateElementDistribution(chinese.ElementBalance)
	patternBits := uint16(variance * 7)
	bits |= patternBits

	// Convert to 4-char hex
	return fmt.Sprintf("%04x", bits)
}

// compressFusionProfile compresses Fusion profile to 4 hex chars (16 bits)
func compressFusionProfile(fusion *FusionProfile) string {
	// Allocate 16 bits:
	// - 4 bits: cognitive fusion pattern
	// - 3 bits: tempo pace
	// - 3 bits: intensity
	// - 3 bits: unified balance
	// - 3 bits: harmonic resonance

	var bits uint16

	// Cognitive fusion pattern (bits 15-12)
	// Encode which cognitive aspect is dominant
	cogPattern := getCognitivePattern(fusion.CognitiveFusion)
	bits |= (cogPattern << 12)

	// Tempo pace (bits 11-9)
	paceBits := uint16(fusion.TempoSignals.Pace * 7)
	bits |= (paceBits << 9)

	// Intensity (bits 8-6)
	intensityBits := uint16(fusion.TempoSignals.Intensity * 7)
	bits |= (intensityBits << 6)

	// Unified balance (bits 5-3)
	balanceBits := uint16(fusion.UnifiedBalance * 7)
	bits |= (balanceBits << 3)

	// Harmonic resonance (bits 2-0)
	resonanceBits := uint16(fusion.HarmonicResonance * 7)
	bits |= resonanceBits

	// Convert to 4-char hex
	return fmt.Sprintf("%04x", bits)
}

// generateU5Chip generates a unique CHIP for U5 using all profile data
func generateU5Chip(western *WesternProfile, chinese *ChineseProfile, fusion *FusionProfile, salt []byte) (string, error) {
	// Create a deterministic string representation of all profiles
	data := fmt.Sprintf("U5|W:%+v|C:%+v|F:%+v", western, chinese, fusion)

	// Concatenate salt + data
	input := append(salt, []byte(data)...)

	// Compute SHA256
	hash := sha256.Sum256(input)

	// Take first 12 hex characters (48 bits)
	hexHash := hex.EncodeToString(hash[:])
	chip := hexHash[:12]

	return chip, nil
}

// Helper function to calculate element distribution variance
func calculateElementDistribution(elements map[string]float64) float64 {
	// Calculate how evenly distributed the elements are
	mean := 0.2 // Expected mean for 5 elements
	sumSquaredDiff := 0.0

	for _, value := range elements {
		diff := value - mean
		sumSquaredDiff += diff * diff
	}

	// Return normalized variance (0 = perfectly even, 1 = highly skewed)
	variance := sumSquaredDiff / 5
	return clampValue(variance * 10) // Scale and clamp
}

// getCognitivePattern returns a 4-bit pattern representing dominant cognitive trait
func getCognitivePattern(cog CognitiveFusion) uint16 {
	maxVal := 0.0
	pattern := uint16(0)

	traits := map[uint16]float64{
		0: cog.Analytical,
		1: cog.Creative,
		2: cog.Grounded,
		3: cog.Adaptive,
		4: cog.Expressive,
	}

	for p, val := range traits {
		if val > maxVal {
			maxVal = val
			pattern = p
		}
	}

	// Add secondary trait influence (bits 3-2)
	secondMax := 0.0
	secondPattern := uint16(0)
	for p, val := range traits {
		if p != pattern && val > secondMax {
			secondMax = val
			secondPattern = p
		}
	}

	// Combine primary and secondary patterns
	return (pattern << 2) | (secondPattern & 0x3)
}

// ValidateU5Format checks if a string matches the expected HCS-U5 format
func ValidateU5Format(code string) bool {
	// Expected format: HCS-U5|XX|W:xxxx|C:xxxx|F:xxxx|CHIP:xxxxxxxxxxxx
	if len(code) < 30 { // Minimum viable length
		return false
	}

	// Check prefix
	if len(code) < 7 || code[:7] != "HCS-U5|" {
		return false
	}

	// Basic structure validation (can be enhanced with regex)
	// Check for presence of required segments
	requiredSegments := []string{"|W:", "|C:", "|F:", "|CHIP:"}
	for _, seg := range requiredSegments {
		found := false
		for i := 0; i < len(code)-len(seg); i++ {
			if code[i:i+len(seg)] == seg {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// DecodeU5 decodes an HCS-U5 code to extract basic information
// This is a simplified decoder for verification purposes
func DecodeU5(code string) (map[string]string, error) {
	if !ValidateU5Format(code) {
		return nil, fmt.Errorf("invalid HCS-U5 format")
	}

	components := make(map[string]string)

	// Extract fusion ID (chars 7-8)
	components["fusionId"] = code[7:9]

	// Extract segments using simple parsing
	// Find positions of each segment
	wPos := -1
	cPos := -1
	fPos := -1
	chipPos := -1

	for i := 0; i < len(code)-2; i++ {
		if code[i:i+2] == "W:" && wPos == -1 {
			wPos = i + 2
		} else if code[i:i+2] == "C:" && cPos == -1 {
			cPos = i + 2
		} else if code[i:i+2] == "F:" && fPos == -1 {
			fPos = i + 2
		}
	}

	for i := 0; i < len(code)-5; i++ {
		if code[i:i+5] == "CHIP:" {
			chipPos = i + 5
			break
		}
	}

	// Extract hex values
	if wPos > 0 && cPos > wPos {
		components["western"] = code[wPos : wPos+4]
	}
	if cPos > 0 && fPos > cPos {
		components["chinese"] = code[cPos : cPos+4]
	}
	if fPos > 0 && chipPos > fPos {
		components["fusion"] = code[fPos : fPos+4]
	}
	if chipPos > 0 && chipPos+12 <= len(code) {
		components["chip"] = code[chipPos : chipPos+12]
	}

	return components, nil
}
