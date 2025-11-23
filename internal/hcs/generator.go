package hcs

import (
	"fmt"
)

// Generator handles HCS code generation with persistent salt
type Generator struct {
	salt []byte
}

// GeneratorOptions allows customization of code generation
type GeneratorOptions struct {
	U3Only bool // Only generate U3 code
	U4Only bool // Only generate U4 code
}

// NewGenerator creates a new HCS code generator
func NewGenerator() (*Generator, error) {
	// Load or create salt from current directory
	salt, err := LoadOrCreateSalt(".")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize generator: %w", err)
	}

	return &Generator{
		salt: salt,
	}, nil
}

// NewGeneratorWithSaltDir creates a generator with a specific salt directory
func NewGeneratorWithSaltDir(dir string) (*Generator, error) {
	salt, err := LoadOrCreateSalt(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize generator with dir %s: %w", dir, err)
	}

	return &Generator{
		salt: salt,
	}, nil
}

// Generate creates HCS codes from an input profile
func (g *Generator) Generate(in *InputProfile) (*OutputHCS, error) {
	return g.GenerateWithOptions(in, nil)
}

// GenerateWithOptions creates HCS codes with specific options
func (g *Generator) GenerateWithOptions(in *InputProfile, opts *GeneratorOptions) (*OutputHCS, error) {
	if in == nil {
		return nil, fmt.Errorf("input profile cannot be nil")
	}

	// Validate input
	if err := g.validateInput(in); err != nil {
		return nil, fmt.Errorf("invalid input profile: %w", err)
	}

	// Default options
	if opts == nil {
		opts = &GeneratorOptions{}
	}

	// Normalize the profile for consistent processing
	normalized := NormalizeProfile(in)

	// Generate CHIP signature
	chip, err := GenerateCHIP(g.salt, normalized)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CHIP: %w", err)
	}

	output := &OutputHCS{
		Input: *in,
		Chip:  chip,
	}

	// Generate U3 code unless U4Only is set
	if !opts.U4Only {
		output.CodeU3 = EncodeU3(in, chip)
	}

	// Generate U4 code unless U3Only is set
	if !opts.U3Only {
		u4Code, err := EncodeU4(normalized, chip)
		if err != nil {
			return nil, fmt.Errorf("failed to generate U4 code: %w", err)
		}
		output.CodeU4 = u4Code
	}

	// Generate Chinese profile and U5 if birth info is provided
	if in.BirthInfo != nil {
		// Compute Chinese BaZi profile
		chineseProfile, err := ComputeChineseProfile(*in.BirthInfo)
		if err != nil {
			// Log error but don't fail the entire generation
			// Chinese profile is optional enhancement
			fmt.Printf("Warning: failed to compute Chinese profile: %v\n", err)
		} else {
			output.ChineseProfile = chineseProfile

			// Create Western profile from input
			westernProfile := &WesternProfile{
				DominantElement: in.DominantElement,
				Modal:           in.Modal,
				Cognition:       in.Cognition,
				Interaction:     in.Interaction,
			}

			// Build fusion profile
			fusionProfile := BuildFusionProfile(westernProfile, chineseProfile)

			// Create combined profile
			output.CombinedProfile = &CombinedProfile{
				Western: *westernProfile,
				Chinese: *chineseProfile,
				Fusion:  *fusionProfile,
			}

			// Generate HCS-U5 code
			u5Code, err := EncodeU5(westernProfile, chineseProfile, fusionProfile, g.salt)
			if err != nil {
				fmt.Printf("Warning: failed to generate U5 code: %v\n", err)
			} else {
				output.CodeU5 = u5Code
			}
		}
	}

	// Generate canonical profile data for U7 signatures (uses normalized + optional combined profile)
	canonical, err := CanonicalProfileData(normalized, output.CombinedProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to build canonical profile: %w", err)
	}

	// Load secret key from environment (required for U7). This is a hard failure
	// to avoid accidentally generating unsigned or weakly signed codes.
	secret, err := LoadSecretKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load secret key: %w", err)
	}

	// Compute quantum-style signatures using the canonical data, secret key, and persistent salt.
	qsigHex, b3Hex, err := ComputeQuantumSignatures(canonical, secret, g.salt)
	if err != nil {
		return nil, fmt.Errorf("failed to compute quantum signatures: %w", err)
	}

	// Format HCS-U7 code using the normalized profile and signatures.
	u7, err := FormatHCSU7(normalized, qsigHex, b3Hex)
	if err != nil {
		return nil, fmt.Errorf("failed to format HCS-U7 code: %w", err)
	}

	output.CodeU7 = u7
	output.QSig = qsigHex
	output.B3Sig = b3Hex

	return output, nil
}

// validateInput checks if the input profile has valid values
func (g *Generator) validateInput(in *InputProfile) error {
	// Validate element
	validElements := map[string]bool{
		"Earth": true,
		"Air":   true,
		"Water": true,
		"Fire":  true,
	}
	if !validElements[in.DominantElement] {
		return fmt.Errorf("invalid dominant element: %s", in.DominantElement)
	}

	// Validate modal values (should be between 0 and 1)
	if err := g.validateRange("modal.cardinal", in.Modal.Cardinal); err != nil {
		return err
	}
	if err := g.validateRange("modal.fixed", in.Modal.Fixed); err != nil {
		return err
	}
	if err := g.validateRange("modal.mutable", in.Modal.Mutable); err != nil {
		return err
	}

	// Validate cognition values
	if err := g.validateRange("cognition.fluid", in.Cognition.Fluid); err != nil {
		return err
	}
	if err := g.validateRange("cognition.crystallized", in.Cognition.Crystallized); err != nil {
		return err
	}
	if err := g.validateRange("cognition.verbal", in.Cognition.Verbal); err != nil {
		return err
	}
	if err := g.validateRange("cognition.strategic", in.Cognition.Strategic); err != nil {
		return err
	}
	if err := g.validateRange("cognition.creative", in.Cognition.Creative); err != nil {
		return err
	}

	// Apply sensible defaults for interaction preferences when fields are empty
	if in.Interaction.Pace == "" {
		in.Interaction.Pace = "balanced"
	}
	if in.Interaction.Structure == "" {
		in.Interaction.Structure = "medium"
	}
	if in.Interaction.Tone == "" {
		in.Interaction.Tone = "neutral"
	}

	// Validate interaction preferences
	validPace := map[string]bool{"balanced": true, "fast": true, "slow": true}
	if !validPace[in.Interaction.Pace] {
		return fmt.Errorf("invalid pace: %s", in.Interaction.Pace)
	}

	validStructure := map[string]bool{"low": true, "medium": true, "high": true}
	if !validStructure[in.Interaction.Structure] {
		return fmt.Errorf("invalid structure: %s", in.Interaction.Structure)
	}

	validTone := map[string]bool{"warm": true, "neutral": true, "sharp": true, "precise": true}
	if !validTone[in.Interaction.Tone] {
		return fmt.Errorf("invalid tone: %s", in.Interaction.Tone)
	}

	// Validate optional birth info if provided
	if in.BirthInfo != nil {
		if err := validateBirthInfo(*in.BirthInfo); err != nil {
			return fmt.Errorf("invalid birth info: %w", err)
		}
	}

	return nil
}

// validateRange checks if a value is between 0 and 1
func (g *Generator) validateRange(field string, value float64) error {
	if value < 0 || value > 1 {
		return fmt.Errorf("%s must be between 0 and 1, got %f", field, value)
	}
	return nil
}

// GetSalt returns the current salt (for testing purposes)
func (g *Generator) GetSalt() []byte {
	return g.salt
}
