package hcs

import (
	"encoding/json"
	"fmt"
	"sort"
)

type fixed4 float64

func (f fixed4) MarshalJSON() ([]byte, error) {
	// Always encode with exactly 4 decimal places for canonical stability
	return []byte(fmt.Sprintf("%.4f", float64(f))), nil
}

type canonicalChineseElement struct {
	Name  string `json:"name"`
	Value fixed4 `json:"value"`
}

type canonicalChinese struct {
	YearPillar        string                    `json:"yearPillar"`
	MonthPillar       string                    `json:"monthPillar"`
	DayPillar         string                    `json:"dayPillar"`
	HourPillar        string                    `json:"hourPillar"`
	YinYangBalance    fixed4                    `json:"yinYangBalance"`
	ElementBalance    []canonicalChineseElement `json:"elementBalance,omitempty"`
	DayMaster         string                    `json:"dayMaster"`
	DayMasterStrength fixed4                    `json:"dayMasterStrength"`
}

type canonicalFusion struct {
	FusionID          string `json:"fusionId"`
	UnifiedBalance    fixed4 `json:"unifiedBalance"`
	HarmonicResonance fixed4 `json:"harmonicResonance"`
}

type canonicalProfile struct {
	Normalized *NormalizedProfile `json:"normalized"`
	Chinese    *canonicalChinese  `json:"chinese,omitempty"`
	Fusion     *canonicalFusion   `json:"fusion,omitempty"`
}

// CanonicalProfileData builds a canonical JSON representation of the profile data
// used for cryptographic signatures. The output is deterministic across runs.
func CanonicalProfileData(normalized *NormalizedProfile, combined *CombinedProfile) ([]byte, error) {
	cp := canonicalProfile{
		Normalized: normalized,
	}

	if combined != nil {
		// Chinese component
		ch := &combined.Chinese
		var elems []canonicalChineseElement
		if len(ch.ElementBalance) > 0 {
			keys := make([]string, 0, len(ch.ElementBalance))
			for k := range ch.ElementBalance {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := ch.ElementBalance[k]
				elems = append(elems, canonicalChineseElement{
					Name:  k,
					Value: fixed4(v),
				})
			}
		}

		cp.Chinese = &canonicalChinese{
			YearPillar:        ch.YearPillar,
			MonthPillar:       ch.MonthPillar,
			DayPillar:         ch.DayPillar,
			HourPillar:        ch.HourPillar,
			YinYangBalance:    fixed4(ch.YinYangBalance),
			ElementBalance:    elems,
			DayMaster:         ch.DayMaster,
			DayMasterStrength: fixed4(ch.DayMasterStrength),
		}

		// Fusion component
		f := combined.Fusion
		cp.Fusion = &canonicalFusion{
			FusionID:          f.FusionID,
			UnifiedBalance:    fixed4(f.UnifiedBalance),
			HarmonicResonance: fixed4(f.HarmonicResonance),
		}
	}

	b, err := json.Marshal(cp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal canonical profile: %w", err)
	}
	return b, nil
}
