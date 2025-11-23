package tests

import (
	"testing"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

// TestComputeYearPillar tests the Year Pillar computation
func TestComputeYearPillar(t *testing.T) {
	testCases := []struct {
		year     int
		expected string
	}{
		{1984, "Jia-Zi"},    // Wood Rat year
		{1990, "Geng-Wu"},   // Metal Horse year
		{2000, "Geng-Chen"}, // Metal Dragon year
		{2024, "Jia-Chen"},  // Wood Dragon year
	}

	for _, tc := range testCases {
		t.Run(string(rune(tc.year)), func(t *testing.T) {
			pillar := hcs.ComputeYearPillar(tc.year)
			result := pillar.PillarToString()
			if result != tc.expected {
				t.Errorf("Year %d: expected %s, got %s", tc.year, tc.expected, result)
			}
		})
	}
}

// TestComputeDayPillar tests Day Pillar computation
func TestComputeDayPillar(t *testing.T) {
	// Test deterministic day pillar calculation
	pillar1 := hcs.ComputeDayPillar(1990, 5, 15)
	pillar2 := hcs.ComputeDayPillar(1990, 5, 15)

	if pillar1.PillarToString() != pillar2.PillarToString() {
		t.Error("Day pillar calculation is not deterministic")
	}

	// Different days should produce different pillars (usually)
	pillar3 := hcs.ComputeDayPillar(1990, 5, 16)
	if pillar1.PillarToString() == pillar3.PillarToString() {
		t.Log("Warning: consecutive days produced same pillar (can happen in 60-day cycle)")
	}
}

// TestComputeChineseProfile tests the complete Chinese profile generation
func TestComputeChineseProfile(t *testing.T) {
	birthInfo := hcs.BirthInfo{
		Year:     1990,
		Month:    6,
		Day:      15,
		Hour:     14,
		Minute:   30,
		Timezone: "UTC",
	}

	profile, err := hcs.ComputeChineseProfile(birthInfo)
	if err != nil {
		t.Fatalf("Failed to compute Chinese profile: %v", err)
	}

	// Check that all required fields are present
	if profile.YearPillar == "" {
		t.Error("Year pillar is empty")
	}
	if profile.MonthPillar == "" {
		t.Error("Month pillar is empty")
	}
	if profile.DayPillar == "" {
		t.Error("Day pillar is empty")
	}
	if profile.HourPillar == "" {
		t.Error("Hour pillar is empty")
	}
	if profile.DayMaster == "" {
		t.Error("Day master is empty")
	}

	// Check Yin/Yang balance is within valid range
	if profile.YinYangBalance < 0 || profile.YinYangBalance > 1 {
		t.Errorf("Invalid Yin/Yang balance: %f", profile.YinYangBalance)
	}

	// Check element balance adds up to 1.0 (allowing small floating point error)
	totalElements := 0.0
	for _, val := range profile.ElementBalance {
		totalElements += val
	}
	if totalElements < 0.99 || totalElements > 1.01 {
		t.Errorf("Element balance doesn't sum to 1.0: %f", totalElements)
	}

	// Check all 5 elements are present
	expectedElements := []string{"Wood", "Fire", "Earth", "Metal", "Water"}
	for _, elem := range expectedElements {
		if _, exists := profile.ElementBalance[elem]; !exists {
			t.Errorf("Missing element in balance: %s", elem)
		}
	}
}

// TestElementBalance tests element balance calculation
func TestElementBalance(t *testing.T) {
	// Create sample pillars
	pillars := []hcs.Pillar{
		{Stem: "Jia", Branch: "Zi", StemIndex: 0, BranchIndex: 0},    // Wood/Water
		{Stem: "Bing", Branch: "Wu", StemIndex: 2, BranchIndex: 6},   // Fire/Fire
		{Stem: "Wu", Branch: "Chen", StemIndex: 4, BranchIndex: 4},   // Earth/Earth
		{Stem: "Geng", Branch: "Shen", StemIndex: 6, BranchIndex: 8}, // Metal/Metal
	}

	balance := hcs.CalculateElementBalance(pillars)

	// Check all elements are present
	for _, elem := range []string{"Wood", "Fire", "Earth", "Metal", "Water"} {
		if _, exists := balance[elem]; !exists {
			t.Errorf("Missing element %s in balance", elem)
		}
	}

	// Check total is 1.0
	total := 0.0
	for _, val := range balance {
		total += val
	}
	if total < 0.99 || total > 1.01 {
		t.Errorf("Element balance doesn't sum to 1.0: %f", total)
	}
}

// TestYinYangBalance tests Yin/Yang balance calculation
func TestYinYangBalance(t *testing.T) {
	// All Yang pillars
	yangPillars := []hcs.Pillar{
		{Stem: "Jia", Branch: "Zi", StemIndex: 0, BranchIndex: 0},    // Yang/Yang
		{Stem: "Bing", Branch: "Wu", StemIndex: 2, BranchIndex: 6},   // Yang/Yang
		{Stem: "Wu", Branch: "Chen", StemIndex: 4, BranchIndex: 4},   // Yang/Yang
		{Stem: "Geng", Branch: "Shen", StemIndex: 6, BranchIndex: 8}, // Yang/Yang
	}

	yangBalance := hcs.CalculateYinYangBalance(yangPillars)
	if yangBalance < 0.9 {
		t.Errorf("Expected high Yang balance for all Yang pillars, got %f", yangBalance)
	}

	// All Yin pillars
	yinPillars := []hcs.Pillar{
		{Stem: "Yi", Branch: "Chou", StemIndex: 1, BranchIndex: 1}, // Yin/Yin
		{Stem: "Ding", Branch: "Si", StemIndex: 3, BranchIndex: 5}, // Yin/Yin
		{Stem: "Ji", Branch: "Wei", StemIndex: 5, BranchIndex: 7},  // Yin/Yin
		{Stem: "Xin", Branch: "You", StemIndex: 7, BranchIndex: 9}, // Yin/Yin
	}

	yinBalance := hcs.CalculateYinYangBalance(yinPillars)
	if yinBalance > 0.1 {
		t.Errorf("Expected low Yang balance (high Yin) for all Yin pillars, got %f", yinBalance)
	}
}

// TestValidateBirthInfo tests birth info validation
func TestValidateBirthInfo(t *testing.T) {
	validInfo := hcs.BirthInfo{
		Year:     1990,
		Month:    6,
		Day:      15,
		Hour:     14,
		Minute:   30,
		Timezone: "UTC",
	}

	// Should be valid
	_, err := hcs.ComputeChineseProfile(validInfo)
	if err != nil {
		t.Errorf("Valid birth info rejected: %v", err)
	}

	// Invalid year
	invalidYear := validInfo
	invalidYear.Year = 1800
	_, err = hcs.ComputeChineseProfile(invalidYear)
	if err == nil {
		t.Error("Invalid year (1800) was not rejected")
	}

	// Invalid month
	invalidMonth := validInfo
	invalidMonth.Month = 13
	_, err = hcs.ComputeChineseProfile(invalidMonth)
	if err == nil {
		t.Error("Invalid month (13) was not rejected")
	}

	// Invalid hour
	invalidHour := validInfo
	invalidHour.Hour = 25
	_, err = hcs.ComputeChineseProfile(invalidHour)
	if err == nil {
		t.Error("Invalid hour (25) was not rejected")
	}
}
