package hcs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
)

// NormalizedProfile represents the normalized values for canonical hashing
type NormalizedProfile struct {
	Element string                `json:"element"`
	Modal   NormalizedModal       `json:"modal"`
	Cog     NormalizedCognition   `json:"cog"`
	Int     NormalizedInteraction `json:"int"`
}

// NormalizedModal with integer percentages
type NormalizedModal struct {
	C int `json:"c"` // Cardinal
	F int `json:"f"` // Fixed
	M int `json:"m"` // Mutable
}

// NormalizedCognition with integer percentages
type NormalizedCognition struct {
	F  int `json:"F"`  // Fluid
	C  int `json:"C"`  // Crystallized
	V  int `json:"V"`  // Verbal
	S  int `json:"S"`  // Strategic
	Cr int `json:"Cr"` // Creative
}

// NormalizedInteraction with single letter codes
type NormalizedInteraction struct {
	PB string `json:"PB"` // Pace
	SM string `json:"SM"` // Structure
	TN string `json:"TN"` // Tone
}

// NormalizeProfile converts an InputProfile to normalized values for hashing
func NormalizeProfile(in *InputProfile) *NormalizedProfile {
	return &NormalizedProfile{
		Element: mapElementToLetter(in.DominantElement),
		Modal: NormalizedModal{
			C: clampAndRound(in.Modal.Cardinal),
			F: clampAndRound(in.Modal.Fixed),
			M: clampAndRound(in.Modal.Mutable),
		},
		Cog: NormalizedCognition{
			F:  clampAndRound(in.Cognition.Fluid),
			C:  clampAndRound(in.Cognition.Crystallized),
			V:  clampAndRound(in.Cognition.Verbal),
			S:  clampAndRound(in.Cognition.Strategic),
			Cr: clampAndRound(in.Cognition.Creative),
		},
		Int: NormalizedInteraction{
			PB: mapPaceToLetter(in.Interaction.Pace),
			SM: mapStructureToLetter(in.Interaction.Structure),
			TN: mapToneToLetter(in.Interaction.Tone),
		},
	}
}

// GenerateCHIP computes the CHIP-96 (12 hex chars) from salt and normalized profile
func GenerateCHIP(salt []byte, normalized *NormalizedProfile) (string, error) {
	// Create canonical JSON with fixed field order
	canonicalJSON, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("failed to marshal normalized profile: %w", err)
	}

	// Concatenate salt + canonical JSON
	data := append(salt, canonicalJSON...)

	// Compute SHA256
	hash := sha256.Sum256(data)

	// Take first 12 hex characters (48 bits)
	hexHash := hex.EncodeToString(hash[:])
	chip := hexHash[:12]

	return chip, nil
}

// clampAndRound clamps a float between 0 and 1, then converts to percentage
func clampAndRound(value float64) int {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return int(math.Round(value * 100))
}

// mapElementToLetter maps element name to single letter
func mapElementToLetter(element string) string {
	switch element {
	case "Earth":
		return "E"
	case "Air":
		return "A"
	case "Water":
		return "W"
	case "Fire":
		return "F"
	default:
		return "E" // Default to Earth if unknown
	}
}

// mapPaceToLetter maps pace preference to single letter
func mapPaceToLetter(pace string) string {
	switch pace {
	case "balanced":
		return "B"
	case "fast":
		return "F"
	case "slow":
		return "S"
	default:
		return "B"
	}
}

// mapStructureToLetter maps structure preference to single letter
func mapStructureToLetter(structure string) string {
	switch structure {
	case "low":
		return "L"
	case "medium":
		return "M"
	case "high":
		return "H"
	default:
		return "M"
	}
}

// mapToneToLetter maps tone preference to single letter
func mapToneToLetter(tone string) string {
	switch tone {
	case "warm":
		return "W"
	case "neutral":
		return "N"
	case "sharp":
		return "S"
	case "precise":
		return "P"
	default:
		return "N"
	}
}
