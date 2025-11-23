package hcs

import (
	"math"
)

// FusionProfile represents the synthesis of Western and Chinese astrological profiles
type FusionProfile struct {
	// Combined element signature (Western 4 + Chinese 5)
	ElementSignature map[string]float64 `json:"elementSignature"`

	// Weighted cognitive tendencies from both systems
	CognitiveFusion CognitiveFusion `json:"cognitiveFusion"`

	// Tempo and expressiveness signals
	TempoSignals TempoSignals `json:"tempoSignals"`

	// Unified balance metrics
	UnifiedBalance    float64 `json:"unifiedBalance"`    // 0-1, combined Yin/Yang and modal balance
	HarmonicResonance float64 `json:"harmonicResonance"` // 0-1, how well the two systems align

	// Fusion ID for compact encoding
	FusionID string `json:"fusionId"`
}

// CognitiveFusion represents merged cognitive patterns
type CognitiveFusion struct {
	Analytical float64 `json:"analytical"` // Strategic + Metal/Water influence
	Creative   float64 `json:"creative"`   // Creative + Fire/Wood influence
	Grounded   float64 `json:"grounded"`   // Crystallized + Earth influence
	Adaptive   float64 `json:"adaptive"`   // Fluid + changing elements
	Expressive float64 `json:"expressive"` // Verbal + Yang influence
}

// TempoSignals represents timing and rhythm preferences
type TempoSignals struct {
	Pace        float64 `json:"pace"`        // 0 = slow, 0.5 = balanced, 1 = fast
	Variability float64 `json:"variability"` // 0 = consistent, 1 = highly variable
	Intensity   float64 `json:"intensity"`   // 0 = gentle, 1 = intense
	Rhythm      string  `json:"rhythm"`      // "steady", "dynamic", "fluctuating"
}

// CombinedProfile contains Western, Chinese, and Fusion profiles
type CombinedProfile struct {
	Western WesternProfile `json:"western"`
	Chinese ChineseProfile `json:"chinese"`
	Fusion  FusionProfile  `json:"fusion"`
}

// WesternProfile represents the Western astrological profile (based on existing InputProfile)
type WesternProfile struct {
	DominantElement string                 `json:"dominantElement"`
	Modal           ModalBalance           `json:"modal"`
	Cognition       CognitionProfile       `json:"cognition"`
	Interaction     InteractionPreferences `json:"interaction"`
}

// BuildFusionProfile creates a fusion profile from Western and Chinese profiles
func BuildFusionProfile(western *WesternProfile, chinese *ChineseProfile) *FusionProfile {
	// Create element signature combining both systems
	elementSig := buildElementSignature(western, chinese)

	// Build cognitive fusion
	cogFusion := buildCognitiveFusion(western, chinese)

	// Build tempo signals
	tempoSignals := buildTempoSignals(western, chinese)

	// Calculate unified balance
	unifiedBalance := calculateUnifiedBalance(western, chinese)

	// Calculate harmonic resonance (how well systems align)
	harmonicResonance := calculateHarmonicResonance(western, chinese)

	// Generate fusion ID
	fusionID := generateFusionID(western, chinese)

	return &FusionProfile{
		ElementSignature:  elementSig,
		CognitiveFusion:   cogFusion,
		TempoSignals:      tempoSignals,
		UnifiedBalance:    unifiedBalance,
		HarmonicResonance: harmonicResonance,
		FusionID:          fusionID,
	}
}

// buildElementSignature combines Western 4 elements with Chinese 5 elements
func buildElementSignature(western *WesternProfile, chinese *ChineseProfile) map[string]float64 {
	signature := make(map[string]float64)

	// Map Western elements to unified schema
	// Western uses: Earth, Air, Water, Fire
	westernWeight := 0.4 // 40% influence from Western

	switch western.DominantElement {
	case "Earth":
		signature["Earth"] += westernWeight
	case "Air":
		// Air maps to Wood and Metal in Chinese system
		signature["Wood"] += westernWeight * 0.5
		signature["Metal"] += westernWeight * 0.5
	case "Water":
		signature["Water"] += westernWeight
	case "Fire":
		signature["Fire"] += westernWeight
	}

	// Add Chinese elements (60% influence)
	chineseWeight := 0.6
	for element, balance := range chinese.ElementBalance {
		signature[element] += balance * chineseWeight
	}

	// Normalize to sum to 1
	total := 0.0
	for _, v := range signature {
		total += v
	}
	if total > 0 {
		for k := range signature {
			signature[k] /= total
		}
	}

	return signature
}

// buildCognitiveFusion merges cognitive patterns from both systems
func buildCognitiveFusion(western *WesternProfile, chinese *ChineseProfile) CognitiveFusion {
	// Extract Chinese element influences
	metalInfluence := chinese.ElementBalance["Metal"]
	waterInfluence := chinese.ElementBalance["Water"]
	fireInfluence := chinese.ElementBalance["Fire"]
	woodInfluence := chinese.ElementBalance["Wood"]
	earthInfluence := chinese.ElementBalance["Earth"]
	yangInfluence := chinese.YinYangBalance

	// Analytical: Strategic thinking + Metal/Water clarity
	analytical := western.Cognition.Strategic*0.5 +
		(metalInfluence*0.3 + waterInfluence*0.2)

	// Creative: Creative + Fire/Wood growth
	creative := western.Cognition.Creative*0.5 +
		(fireInfluence*0.3 + woodInfluence*0.2)

	// Grounded: Crystallized knowledge + Earth stability
	grounded := western.Cognition.Crystallized*0.5 +
		earthInfluence*0.5

	// Adaptive: Fluid intelligence + element variability
	elementVariability := calculateElementVariability(chinese.ElementBalance)
	adaptive := western.Cognition.Fluid*0.6 + elementVariability*0.4

	// Expressive: Verbal + Yang energy
	expressive := western.Cognition.Verbal*0.5 + yangInfluence*0.5

	return CognitiveFusion{
		Analytical: clampValue(analytical),
		Creative:   clampValue(creative),
		Grounded:   clampValue(grounded),
		Adaptive:   clampValue(adaptive),
		Expressive: clampValue(expressive),
	}
}

// buildTempoSignals creates tempo and rhythm preferences
func buildTempoSignals(western *WesternProfile, chinese *ChineseProfile) TempoSignals {
	// Pace influenced by Western pace and Chinese Yang energy
	basePace := 0.5 // balanced default
	switch western.Interaction.Pace {
	case "fast":
		basePace = 0.8
	case "slow":
		basePace = 0.2
	}
	// Yang energy speeds up, Yin slows down
	pace := basePace*0.6 + chinese.YinYangBalance*0.4

	// Variability from modal balance and element distribution
	modalVariability := calculateModalVariability(western.Modal)
	elementVariability := calculateElementVariability(chinese.ElementBalance)
	variability := (modalVariability + elementVariability) / 2

	// Intensity from Fire/Water balance and Day Master strength
	fireWater := chinese.ElementBalance["Fire"] + chinese.ElementBalance["Water"]
	intensity := fireWater*0.5 + chinese.DayMasterStrength*0.5

	// Determine rhythm pattern
	rhythm := "steady"
	if variability > 0.6 {
		rhythm = "fluctuating"
	} else if intensity > 0.6 && pace > 0.6 {
		rhythm = "dynamic"
	}

	return TempoSignals{
		Pace:        clampValue(pace),
		Variability: clampValue(variability),
		Intensity:   clampValue(intensity),
		Rhythm:      rhythm,
	}
}

// calculateUnifiedBalance combines Yin/Yang with modal balance
func calculateUnifiedBalance(western *WesternProfile, chinese *ChineseProfile) float64 {
	// Modal balance center (Cardinal vs Fixed primarily)
	modalBalance := western.Modal.Cardinal*0.5 + western.Modal.Mutable*0.3 + (1-western.Modal.Fixed)*0.2

	// Combine with Yin/Yang
	unified := modalBalance*0.4 + chinese.YinYangBalance*0.6

	return clampValue(unified)
}

// calculateHarmonicResonance measures how well the two systems align
func calculateHarmonicResonance(western *WesternProfile, chinese *ChineseProfile) float64 {
	resonance := 0.5 // Base resonance

	// Check element compatibility
	westernElement := western.DominantElement
	chineseDominant := getDominantElement(chinese.ElementBalance)

	// Compatible elements increase resonance
	if areElementsCompatible(westernElement, chineseDominant) {
		resonance += 0.2
	}

	// Check pace alignment
	if western.Interaction.Pace == "fast" && chinese.YinYangBalance > 0.6 {
		resonance += 0.15
	} else if western.Interaction.Pace == "slow" && chinese.YinYangBalance < 0.4 {
		resonance += 0.15
	} else if western.Interaction.Pace == "balanced" &&
		chinese.YinYangBalance >= 0.4 && chinese.YinYangBalance <= 0.6 {
		resonance += 0.15
	}

	// Check structure vs element stability
	earthMetal := chinese.ElementBalance["Earth"] + chinese.ElementBalance["Metal"]
	if western.Interaction.Structure == "high" && earthMetal > 0.4 {
		resonance += 0.15
	} else if western.Interaction.Structure == "low" && earthMetal < 0.3 {
		resonance += 0.15
	}

	return clampValue(resonance)
}

// generateFusionID creates a 2-character fusion identifier
func generateFusionID(western *WesternProfile, chinese *ChineseProfile) string {
	// First character based on combined elements
	elementCode := getElementCode(western.DominantElement,
		getDominantElement(chinese.ElementBalance))

	// Second character based on balance and tempo
	balanceCode := getBalanceCode(western.Modal, chinese.YinYangBalance)

	return elementCode + balanceCode
}

// Helper functions

func calculateElementVariability(elements map[string]float64) float64 {
	// Calculate standard deviation of element distribution
	mean := 0.2 // Expected mean for 5 elements
	variance := 0.0

	for _, value := range elements {
		diff := value - mean
		variance += diff * diff
	}

	// Normalize variance to 0-1 range
	stdDev := math.Sqrt(variance / 5)
	return math.Min(stdDev*4, 1.0) // Scale up and cap at 1
}

func calculateModalVariability(modal ModalBalance) float64 {
	// Calculate how distributed the modalities are
	values := []float64{modal.Cardinal, modal.Fixed, modal.Mutable}
	mean := 1.0 / 3.0
	variance := 0.0

	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}

	stdDev := math.Sqrt(variance / 3)
	return 1.0 - math.Min(stdDev*3, 1.0) // Invert so high variability = close to 1
}

func getDominantElement(elements map[string]float64) string {
	maxElement := ""
	maxValue := 0.0

	for element, value := range elements {
		if value > maxValue {
			maxValue = value
			maxElement = element
		}
	}

	return maxElement
}

func areElementsCompatible(western, chinese string) bool {
	// Define compatibility rules
	compatible := map[string][]string{
		"Fire":  {"Fire", "Wood"},
		"Earth": {"Earth", "Fire", "Metal"},
		"Air":   {"Wood", "Metal"},
		"Water": {"Water", "Wood"},
	}

	if compatList, exists := compatible[western]; exists {
		for _, elem := range compatList {
			if elem == chinese {
				return true
			}
		}
	}

	return false
}

func getElementCode(western, chinese string) string {
	// Create a deterministic 1-char code from element combination
	codes := map[string]string{
		"Fire-Fire": "A", "Fire-Wood": "B", "Fire-Earth": "C", "Fire-Metal": "D", "Fire-Water": "E",
		"Earth-Fire": "F", "Earth-Wood": "G", "Earth-Earth": "H", "Earth-Metal": "I", "Earth-Water": "J",
		"Air-Fire": "K", "Air-Wood": "L", "Air-Earth": "M", "Air-Metal": "N", "Air-Water": "O",
		"Water-Fire": "P", "Water-Wood": "Q", "Water-Earth": "R", "Water-Metal": "S", "Water-Water": "T",
	}

	key := western + "-" + chinese
	if code, exists := codes[key]; exists {
		return code
	}
	return "X" // Default
}

func getBalanceCode(modal ModalBalance, yinYang float64) string {
	// Create a deterministic 1-char code from balance metrics
	// Divide into 9 zones (3x3 grid)
	modalZone := 0
	if modal.Cardinal > 0.4 {
		modalZone = 2
	} else if modal.Fixed > 0.4 {
		modalZone = 1
	}

	yinYangZone := 0
	if yinYang > 0.66 {
		yinYangZone = 2
	} else if yinYang > 0.33 {
		yinYangZone = 1
	}

	codes := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	index := modalZone*3 + yinYangZone

	return codes[index]
}

func clampValue(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}
