package hcs

// ModalBalance represents the distribution of astrological modalities
type ModalBalance struct {
	Cardinal float64 `json:"cardinal"`
	Fixed    float64 `json:"fixed"`
	Mutable  float64 `json:"mutable"`
}

// CognitionProfile represents cognitive processing capabilities
type CognitionProfile struct {
	Fluid        float64 `json:"fluid"`
	Crystallized float64 `json:"crystallized"`
	Verbal       float64 `json:"verbal"`
	Strategic    float64 `json:"strategic"`
	Creative     float64 `json:"creative"`
}

// InteractionPreferences represents communication and interaction styles
type InteractionPreferences struct {
	Pace      string `json:"pace"`      // "balanced" | "fast" | "slow"
	Structure string `json:"structure"` // "low" | "medium" | "high"
	Tone      string `json:"tone"`      // "warm" | "neutral" | "sharp" | "precise"
}

// InputProfile represents the complete input for HCS code generation
type InputProfile struct {
	DominantElement string                 `json:"dominantElement"` // "Earth" | "Air" | "Water" | "Fire"
	Modal           ModalBalance           `json:"modal"`
	Cognition       CognitionProfile       `json:"cognition"`
	Interaction     InteractionPreferences `json:"interaction"`
	// Optional birth info for Chinese astrology
	BirthInfo *BirthInfo `json:"birthInfo,omitempty"`
}

// OutputHCS represents the generated HCS codes and metadata
type OutputHCS struct {
	Input           InputProfile     `json:"input"`
	CodeU3          string           `json:"codeU3"`
	CodeU4          string           `json:"codeU4,omitempty"`
	CodeU5          string           `json:"codeU5,omitempty"` // NEW: HCS-U5 fusion code
	Chip            string           `json:"chip"`
	ChineseProfile  *ChineseProfile  `json:"chineseProfile,omitempty"`  // NEW: Chinese BaZi profile
	CombinedProfile *CombinedProfile `json:"combinedProfile,omitempty"` // NEW: Combined profiles
}
