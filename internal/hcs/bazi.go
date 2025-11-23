package hcs

import (
	"math"
	"time"
)

// BaZi constants and tables for Chinese astrology computation

// HeavenlyStems represents the 10 Heavenly Stems
var HeavenlyStems = []struct {
	Name    string
	Element string
	YinYang string
}{
	{"Jia", "Wood", "Yang"},
	{"Yi", "Wood", "Yin"},
	{"Bing", "Fire", "Yang"},
	{"Ding", "Fire", "Yin"},
	{"Wu", "Earth", "Yang"},
	{"Ji", "Earth", "Yin"},
	{"Geng", "Metal", "Yang"},
	{"Xin", "Metal", "Yin"},
	{"Ren", "Water", "Yang"},
	{"Gui", "Water", "Yin"},
}

// EarthlyBranches represents the 12 Earthly Branches
var EarthlyBranches = []struct {
	Name    string
	Element string
	YinYang string
	Animal  string
}{
	{"Zi", "Water", "Yang", "Rat"},
	{"Chou", "Earth", "Yin", "Ox"},
	{"Yin", "Wood", "Yang", "Tiger"},
	{"Mao", "Wood", "Yin", "Rabbit"},
	{"Chen", "Earth", "Yang", "Dragon"},
	{"Si", "Fire", "Yin", "Snake"},
	{"Wu", "Fire", "Yang", "Horse"},
	{"Wei", "Earth", "Yin", "Goat"},
	{"Shen", "Metal", "Yang", "Monkey"},
	{"You", "Metal", "Yin", "Rooster"},
	{"Xu", "Earth", "Yang", "Dog"},
	{"Hai", "Water", "Yin", "Pig"},
}

// MonthBranchMapping maps month numbers to earthly branches
// Based on solar calendar approximation
var MonthBranchMapping = []int{
	2,  // January - Chou
	3,  // February - Yin
	4,  // March - Mao
	5,  // April - Chen
	6,  // May - Si
	7,  // June - Wu
	8,  // July - Wei
	9,  // August - Shen
	10, // September - You
	11, // October - Xu
	0,  // November - Hai
	1,  // December - Zi
}

// Pillar represents a BaZi pillar with Heavenly Stem and Earthly Branch
type Pillar struct {
	Stem       string `json:"stem"`
	Branch     string `json:"branch"`
	StemIndex  int    `json:"stemIndex"`
	BranchIndex int   `json:"branchIndex"`
}

// PillarToString returns the string representation of a pillar
func (p Pillar) PillarToString() string {
	return p.Stem + "-" + p.Branch
}

// GetElement returns the combined element influence of the pillar
func (p Pillar) GetElement() string {
	// Primary element comes from the stem
	return HeavenlyStems[p.StemIndex].Element
}

// GetYinYang returns the Yin/Yang polarity of the pillar
func (p Pillar) GetYinYang() string {
	return HeavenlyStems[p.StemIndex].YinYang
}

// ComputeYearPillar computes the Year Pillar based on birth year
func ComputeYearPillar(year int) Pillar {
	// BaZi year starts from Feb 4 (approximate)
	// Using simple 60-year cycle calculation
	cycleYear := year - 1924 // 1924 is Jia-Zi year (start of cycle)
	
	stemIndex := cycleYear % 10
	if stemIndex < 0 {
		stemIndex += 10
	}
	
	branchIndex := cycleYear % 12
	if branchIndex < 0 {
		branchIndex += 12
	}
	
	return Pillar{
		Stem:        HeavenlyStems[stemIndex].Name,
		Branch:      EarthlyBranches[branchIndex].Name,
		StemIndex:   stemIndex,
		BranchIndex: branchIndex,
	}
}

// ComputeMonthPillar computes the Month Pillar
func ComputeMonthPillar(year int, month int, day int) Pillar {
	// Simplified month pillar calculation
	// In real BaZi, this depends on solar terms
	yearPillar := ComputeYearPillar(year)
	
	// Month branch is determined by solar month
	monthBranchIndex := MonthBranchMapping[month-1]
	
	// Month stem calculation based on year stem
	// Formula: (year stem * 2 + month number) % 10
	monthStemIndex := (yearPillar.StemIndex*2 + month) % 10
	
	return Pillar{
		Stem:        HeavenlyStems[monthStemIndex].Name,
		Branch:      EarthlyBranches[monthBranchIndex].Name,
		StemIndex:   monthStemIndex,
		BranchIndex: monthBranchIndex,
	}
}

// ComputeDayPillar computes the Day Pillar using a deterministic algorithm
func ComputeDayPillar(year, month, day int) Pillar {
	// Using a simplified day pillar calculation
	// Based on days since a reference date
	refDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	targetDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	
	daysSinceRef := int(targetDate.Sub(refDate).Hours() / 24)
	
	// 60-day cycle for stems and branches
	stemIndex := daysSinceRef % 10
	if stemIndex < 0 {
		stemIndex += 10
	}
	
	branchIndex := daysSinceRef % 12
	if branchIndex < 0 {
		branchIndex += 12
	}
	
	return Pillar{
		Stem:        HeavenlyStems[stemIndex].Name,
		Branch:      EarthlyBranches[branchIndex].Name,
		StemIndex:   stemIndex,
		BranchIndex: branchIndex,
	}
}

// ComputeHourPillar computes the Hour Pillar
func ComputeHourPillar(dayPillar Pillar, hour int) Pillar {
	// Hour branch is determined by the time
	// 23:00-01:00 = Zi (0), 01:00-03:00 = Chou (1), etc.
	hourBranchIndex := ((hour + 1) / 2) % 12
	
	// Hour stem calculation based on day stem
	// Formula varies based on day stem
	hourStemBase := (dayPillar.StemIndex % 5) * 2
	hourStemIndex := (hourStemBase + hourBranchIndex) % 10
	
	return Pillar{
		Stem:        HeavenlyStems[hourStemIndex].Name,
		Branch:      EarthlyBranches[hourBranchIndex].Name,
		StemIndex:   hourStemIndex,
		BranchIndex: hourBranchIndex,
	}
}

// CalculateElementBalance calculates the balance of five elements
func CalculateElementBalance(pillars []Pillar) map[string]float64 {
	elements := map[string]float64{
		"Wood":  0,
		"Fire":  0,
		"Earth": 0,
		"Metal": 0,
		"Water": 0,
	}
	
	// Count elements from stems and branches
	for _, pillar := range pillars {
		// Stem element (stronger influence)
		stemElement := HeavenlyStems[pillar.StemIndex].Element
		elements[stemElement] += 1.0
		
		// Branch element (lesser influence)
		branchElement := EarthlyBranches[pillar.BranchIndex].Element
		elements[branchElement] += 0.5
	}
	
	// Normalize to percentages
	total := 0.0
	for _, v := range elements {
		total += v
	}
	
	for k := range elements {
		elements[k] = elements[k] / total
	}
	
	return elements
}

// CalculateYinYangBalance calculates the Yin/Yang balance
func CalculateYinYangBalance(pillars []Pillar) float64 {
	yangCount := 0.0
	totalCount := 0.0
	
	for _, pillar := range pillars {
		// Check stem Yin/Yang
		if HeavenlyStems[pillar.StemIndex].YinYang == "Yang" {
			yangCount += 1.0
		}
		totalCount += 1.0
		
		// Check branch Yin/Yang (lesser weight)
		if EarthlyBranches[pillar.BranchIndex].YinYang == "Yang" {
			yangCount += 0.5
		}
		totalCount += 0.5
	}
	
	// Return Yang percentage (0 = pure Yin, 1 = pure Yang)
	return yangCount / totalCount
}

// GetDayMaster returns the Day Master (Day Stem) information
func GetDayMaster(dayPillar Pillar) string {
	return HeavenlyStems[dayPillar.StemIndex].Name
}

// GetDayMasterStrength evaluates the strength of the Day Master
// Returns a value between 0 (weak) and 1 (strong)
func GetDayMasterStrength(pillars []Pillar, dayPillar Pillar) float64 {
	dayElement := HeavenlyStems[dayPillar.StemIndex].Element
	strength := 0.3 // Base strength
	
	// Check support from other pillars
	for i, pillar := range pillars {
		stemElement := HeavenlyStems[pillar.StemIndex].Element
		branchElement := EarthlyBranches[pillar.BranchIndex].Element
		
		// Skip the day pillar itself
		if i == 2 {
			continue
		}
		
		// Same element strengthens
		if stemElement == dayElement {
			strength += 0.15
		}
		if branchElement == dayElement {
			strength += 0.1
		}
		
		// Generating element strengthens (simplified cycle)
		if isGeneratingElement(stemElement, dayElement) {
			strength += 0.1
		}
		if isGeneratingElement(branchElement, dayElement) {
			strength += 0.05
		}
	}
	
	// Cap between 0 and 1
	return math.Min(math.Max(strength, 0), 1)
}

// isGeneratingElement checks if element1 generates element2 in the creation cycle
func isGeneratingElement(element1, element2 string) bool {
	cycle := map[string]string{
		"Wood":  "Fire",
		"Fire":  "Earth",
		"Earth": "Metal",
		"Metal": "Water",
		"Water": "Wood",
	}
	return cycle[element1] == element2
}
