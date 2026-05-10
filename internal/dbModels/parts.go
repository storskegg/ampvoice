package dbModels

import "gorm.io/gorm"

type Parts struct {
    gorm.Model
    Category   string
    Type       string
    Subtype    string
    Brand      string
    Series     string
    PartNumber string
    Value      string
    Rating     string
    CostUnit   float32
    CostMult   float32
}
