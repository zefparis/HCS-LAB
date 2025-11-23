package hcs

import "fmt"

// FormatHCSU7 assembles the HCS-U7 code from the normalized profile and
// cryptographic signatures. It reuses the same segment semantics as U3/U5
// (E, MOD, COG, INT) while adding quantum-style signature fields.
func FormatHCSU7(profile *NormalizedProfile, qsigHex, b3Hex string) (string, error) {
	if profile == nil {
		return "", fmt.Errorf("normalized profile cannot be nil")
	}
	if len(qsigHex) == 0 || len(b3Hex) == 0 {
		return "", fmt.Errorf("signatures must not be empty")
	}

	// Rebuild segments using the already-normalized integer representation.
	// Element
	elemSegment := fmt.Sprintf("E:%s", profile.Element)

	// Modal
	modalSegment := fmt.Sprintf("MOD:c%02df%02dm%02d", profile.Modal.C, profile.Modal.F, profile.Modal.M)

	// Cognition
	cogSegment := fmt.Sprintf("COG:F%02dC%02dV%02dS%02dCr%02d", profile.Cog.F, profile.Cog.C, profile.Cog.V, profile.Cog.S, profile.Cog.Cr)

	// Interaction
	intSegment := fmt.Sprintf("INT:PB=%s,SM=%s,TN=%s", profile.Int.PB, profile.Int.SM, profile.Int.TN)

	// Truncate signatures for the inline code while keeping full values in JSON metadata.
	qsigInline := qsigHex
	if len(qsigInline) > 24 {
		qsigInline = qsigInline[:24]
	}

	b3Inline := b3Hex
	if len(b3Inline) > 32 {
		b3Inline = b3Inline[:32]
	}

	return fmt.Sprintf(
		"HCS-U7|V:7.0|ALG:QS|%s|%s|%s|%s|QSIG:%s|B3:%s",
		elemSegment,
		modalSegment,
		cogSegment,
		intSegment,
		qsigInline,
		b3Inline,
	), nil
}
