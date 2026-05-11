package dbModels

import "gorm.io/gorm"

type CustomerAmplifier struct {
    gorm.Model
    CustomerId  uint // TODO: FK -> Customers
    AmplifierId uint // TODO: FK -> Amplifiers
    Year        string
}
