package hcs

import (
	"fmt"
	"regexp"
)

// EncodeU3 generates the HCS-U3 code from an InputProfile and CHIP
func EncodeU3(in *InputProfile, chip string) string {
	// Build element segment
	elemSegment := fmt.Sprintf("E:%s", mapElementToLetter(in.DominantElement))

	// Build modal segment
	modalSegment := fmt.Sprintf("MOD:c%02df%02dm%02d",
		clampAndRound(in.Modal.Cardinal),
		clampAndRound(in.Modal.Fixed),
		clampAndRound(in.Modal.Mutable))

	// Build cognition segment
	cogSegment := fmt.Sprintf("COG:F%02dC%02dV%02dS%02dCr%02d",
		clampAndRound(in.Cognition.Fluid),
		clampAndRound(in.Cognition.Crystallized),
		clampAndRound(in.Cognition.Verbal),
		clampAndRound(in.Cognition.Strategic),
		clampAndRound(in.Cognition.Creative))

	// Build interaction segment
	intSegment := fmt.Sprintf("INT:PB=%s,SM=%s,TN=%s",
		mapPaceToLetter(in.Interaction.Pace),
		mapStructureToLetter(in.Interaction.Structure),
		mapToneToLetter(in.Interaction.Tone))

	// Build CHIP segment
	chipSegment := fmt.Sprintf("CHIP:%s", chip)

	// Combine all segments
	return fmt.Sprintf("HCS-U3|%s|%s|%s|%s|%s",
		elemSegment, modalSegment, cogSegment, intSegment, chipSegment)
}

// ValidateU3Format checks if a string matches the expected HCS-U3 format
func ValidateU3Format(code string) bool {
	// Regex pattern for HCS-U3 format
	pattern := `^HCS-U3\|E:[AEWF]\|MOD:c\d{2}f\d{2}m\d{2}\|COG:F\d{2}C\d{2}V\d{2}S\d{2}Cr\d{2}\|INT:PB=[BFS],SM=[LMH],TN=[WNSP]\|CHIP:[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}

// ParseU3 parses an HCS-U3 code and extracts components (optional utility)
func ParseU3(code string) (map[string]string, error) {
	if !ValidateU3Format(code) {
		return nil, fmt.Errorf("invalid HCS-U3 format")
	}

	// This is a simplified parser - you can expand it if needed
	components := make(map[string]string)

	// Extract components using regex groups
	pattern := regexp.MustCompile(`HCS-U3\|E:([AEWF])\|MOD:c(\d{2})f(\d{2})m(\d{2})\|COG:F(\d{2})C(\d{2})V(\d{2})S(\d{2})Cr(\d{2})\|INT:PB=([BFS]),SM=([LMH]),TN=([WNSP])\|CHIP:([0-9a-f]{12})`)
	matches := pattern.FindStringSubmatch(code)

	if len(matches) == 14 {
		components["element"] = matches[1]
		components["modal_cardinal"] = matches[2]
		components["modal_fixed"] = matches[3]
		components["modal_mutable"] = matches[4]
		components["cog_fluid"] = matches[5]
		components["cog_crystallized"] = matches[6]
		components["cog_verbal"] = matches[7]
		components["cog_strategic"] = matches[8]
		components["cog_creative"] = matches[9]
		components["int_pace"] = matches[10]
		components["int_structure"] = matches[11]
		components["int_tone"] = matches[12]
		components["chip"] = matches[13]
	}

	return components, nil
}
