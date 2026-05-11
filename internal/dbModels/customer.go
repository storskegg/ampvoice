package dbModels

import "gorm.io/gorm"

type Customer struct {
    gorm.Model
    NameFirst     string
    NameLast      string
    Email         string
    PhoneNumber   string
    PhoneIsMobile bool
    AddressLine1  string
    AddressLine2  string
    AddressCity   string
    AddressState  string
    AddressZip5   string
    AddressZip4   string
}
