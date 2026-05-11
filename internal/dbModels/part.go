package dbModels

import (
	"gorm.io/gorm"
)

type Part struct {
	gorm.Model
	Category       string  `json:"category"`
	Type           string  `json:"type"`
	SubType        string  `json:"subtype"`
	Brand          string  `json:"brand"`
	Series         string  `json:"series,omitempty"`
	PartNumber     string  `json:"partNumber,omitempty"`
	ValueTolerance string  `json:"valueTolerance,omitempty"`
	Rating         string  `json:"rating,omitempty"`
	Value          string  `json:"value,omitempty"`
	CostUnit       float64 `json:"costUnit"`
	CostMult       float64 `json:"costMult"`
}

func SamplePart() *Part {
	return &Part{
		Category:   "passive",
		Type:       "capacitor",
		SubType:    "electrolytic",
		Brand:      "jupiter",
		Series:     "cosmos",
		PartNumber: "123456789",
		Value:      "100 µF",
		Rating:     "100 V",
		CostUnit:   9.49,
		CostMult:   1.75,
	}
}
