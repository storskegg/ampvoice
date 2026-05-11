package dbModels

import "gorm.io/gorm"

type Amplifier struct {
    gorm.Model
    Make        string
    ModelName   string
    Generation  string
    Power       string
    ModelNumber string
}
