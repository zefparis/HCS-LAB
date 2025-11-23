package hcs

import (
	"fmt"
	"time"
)

// ChineseProfile represents the Chinese BaZi astrological profile
type ChineseProfile struct {
	YearPillar        string             `json:"yearPillar"`
	MonthPillar       string             `json:"monthPillar"`
	DayPillar         string             `json:"dayPillar"`
	HourPillar        string             `json:"hourPillar"`
	YinYangBalance    float64            `json:"yinYangBalance"`    // 0 = pure Yin, 1 = pure Yang
	ElementBalance    map[string]float64 `json:"elementBalance"`    // Wood, Fire, Earth, Metal, Water percentages
	DayMaster         string             `json:"dayMaster"`         // Day stem (most important in BaZi)
	DayMasterStrength float64            `json:"dayMasterStrength"` // 0 = weak, 1 = strong
}

// BirthInfo contains the birth date and time information needed for BaZi
type BirthInfo struct {
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	Timezone string `json:"timezone"`
}

// ComputeChineseProfile generates a complete Chinese astrological profile
func ComputeChineseProfile(birthInfo BirthInfo) (*ChineseProfile, error) {
	// Validate input
	if err := validateBirthInfo(birthInfo); err != nil {
		return nil, err
	}

	// Load timezone if specified
	loc := time.UTC
	if birthInfo.Timezone != "" && birthInfo.Timezone != "UTC" {
		parsedLoc, err := time.LoadLocation(birthInfo.Timezone)
		if err == nil {
			loc = parsedLoc
		}
		// If timezone parsing fails, continue with UTC
	}

	// Create birth time in the specified timezone
	birthTime := time.Date(
		birthInfo.Year,
		time.Month(birthInfo.Month),
		birthInfo.Day,
		birthInfo.Hour,
		birthInfo.Minute,
		0, 0, loc,
	)

	// Convert to local time for BaZi calculation
	// BaZi traditionally uses local solar time
	year := birthTime.Year()
	month := int(birthTime.Month())
	day := birthTime.Day()
	hour := birthTime.Hour()

	// Compute the four pillars
	yearPillar := ComputeYearPillar(year)
	monthPillar := ComputeMonthPillar(year, month, day)
	dayPillar := ComputeDayPillar(year, month, day)
	hourPillar := ComputeHourPillar(dayPillar, hour)

	// Collect all pillars
	pillars := []Pillar{yearPillar, monthPillar, dayPillar, hourPillar}

	// Calculate element balance
	elementBalance := CalculateElementBalance(pillars)

	// Calculate Yin/Yang balance
	yinYangBalance := CalculateYinYangBalance(pillars)

	// Get Day Master
	dayMaster := GetDayMaster(dayPillar)

	// Calculate Day Master strength
	dayMasterStrength := GetDayMasterStrength(pillars, dayPillar)

	return &ChineseProfile{
		YearPillar:        yearPillar.PillarToString(),
		MonthPillar:       monthPillar.PillarToString(),
		DayPillar:         dayPillar.PillarToString(),
		HourPillar:        hourPillar.PillarToString(),
		YinYangBalance:    yinYangBalance,
		ElementBalance:    elementBalance,
		DayMaster:         dayMaster,
		DayMasterStrength: dayMasterStrength,
	}, nil
}

// validateBirthInfo validates the birth information
func validateBirthInfo(info BirthInfo) error {
	// Validate year (reasonable range)
	if info.Year < 1900 || info.Year > 2100 {
		return fmt.Errorf("year must be between 1900 and 2100, got %d", info.Year)
	}

	// Validate month
	if info.Month < 1 || info.Month > 12 {
		return fmt.Errorf("month must be between 1 and 12, got %d", info.Month)
	}

	// Validate day (simplified, doesn't account for month lengths)
	if info.Day < 1 || info.Day > 31 {
		return fmt.Errorf("day must be between 1 and 31, got %d", info.Day)
	}

	// Validate hour
	if info.Hour < 0 || info.Hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23, got %d", info.Hour)
	}

	// Validate minute
	if info.Minute < 0 || info.Minute > 59 {
		return fmt.Errorf("minute must be between 0 and 59, got %d", info.Minute)
	}

	return nil
}

// GetDominantChineseElement returns the most prominent element in the profile
func (cp *ChineseProfile) GetDominantChineseElement() string {
	maxElement := ""
	maxValue := 0.0

	for element, value := range cp.ElementBalance {
		if value > maxValue {
			maxValue = value
			maxElement = element
		}
	}

	return maxElement
}

// GetChineseElementStrength returns the strength of a specific element
func (cp *ChineseProfile) GetChineseElementStrength(element string) float64 {
	if val, exists := cp.ElementBalance[element]; exists {
		return val
	}
	return 0.0
}

// GetYinYangType returns a string representation of the Yin/Yang balance
func (cp *ChineseProfile) GetYinYangType() string {
	if cp.YinYangBalance > 0.6 {
		return "Yang-dominant"
	} else if cp.YinYangBalance < 0.4 {
		return "Yin-dominant"
	}
	return "Balanced"
}

// GetDayMasterType returns a categorization of the Day Master strength
func (cp *ChineseProfile) GetDayMasterType() string {
	if cp.DayMasterStrength > 0.7 {
		return "Strong"
	} else if cp.DayMasterStrength < 0.3 {
		return "Weak"
	}
	return "Moderate"
}

// CompressedChineseData represents a compressed version for encoding
type CompressedChineseData struct {
	YinYang   uint8 // 0-255 scale
	Wood      uint8 // Element percentages as 0-255
	Fire      uint8
	Earth     uint8
	Metal     uint8
	Water     uint8
	DayMaster uint8 // Index of day master stem (0-9)
	Strength  uint8 // Day master strength 0-255
}

// CompressChineseProfile compresses the Chinese profile for compact encoding
func CompressChineseProfile(cp *ChineseProfile) CompressedChineseData {
	// Find Day Master index
	dayMasterIndex := uint8(0)
	for i, stem := range HeavenlyStems {
		if stem.Name == cp.DayMaster {
			dayMasterIndex = uint8(i)
			break
		}
	}

	return CompressedChineseData{
		YinYang:   uint8(cp.YinYangBalance * 255),
		Wood:      uint8(cp.ElementBalance["Wood"] * 255),
		Fire:      uint8(cp.ElementBalance["Fire"] * 255),
		Earth:     uint8(cp.ElementBalance["Earth"] * 255),
		Metal:     uint8(cp.ElementBalance["Metal"] * 255),
		Water:     uint8(cp.ElementBalance["Water"] * 255),
		DayMaster: dayMasterIndex,
		Strength:  uint8(cp.DayMasterStrength * 255),
	}
}
