package dbModels

import "gorm.io/gorm"

type AmplifierPart struct {
    gorm.Model
    AmplifierID uint // TODO: FK -> Amplifiers
    PartId      uint // TODO: FK -> Parts
    Qty         uint
}
