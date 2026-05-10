package dbModels

import (
    "time"

    "gorm.io/gorm"
)

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

type Amplifiers struct {
    gorm.Model
    Make        string
    Model       string
    Generation  string
    Power       string
    ModelNumber string
}

type AmplifierParts struct {
    gorm.Model
    AmplifierID uint // TODO: FK -> Amplifiers
    PartId      uint // TODO: FK -> Parts
    Qty         uint
}

type Customers struct {
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

type CustomerAmplifiers struct {
    gorm.Model
    CustomerId  uint // TODO: FK -> Customers
    AmplifierId uint // TODO: FK -> Amplifiers
    Year        string
}

type ServiceLogs struct {
    CustomerAmplifierId uint // TODO: FK -> CustomerAmplifiers
    TimeStart           time.Time
    TimeEnd             time.Time
    Description         string
    Notes               string
}

type ServiceLogParts struct {
    ServiceLogId uint // TODO: FK -> ServiceLogs
    PartId       uint // TODO: FK -> Parts
    Qty          uint
}

type Invoices struct {
    gorm.Model
}

type InvoiceServiceLogs struct {
    gorm.Model
    InvoiceId    uint // TODO: FK -> Invoices
    ServiceLogId uint // TODO: FK -> ServiceLogs
}

func SamplePart() *Parts {
    return &Parts{
        Category:   "capacitor",
        Type:       "electrolytic",
        Brand:      "jupiter",
        Series:     "cosmos",
        PartNumber: "123456789",
        Value:      "100 uF",
        Rating:     "100 V",
        CostUnit:   9.49,
        CostMult:   1.75,
    }
}
