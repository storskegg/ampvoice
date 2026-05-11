package dbModels

import "gorm.io/gorm"

func MigrateModels(db *gorm.DB) error {
    return db.AutoMigrate(
        &Amplifier{},
        &AmplifierPart{},
        &Customer{},
        &CustomerAmplifier{},
        &Invoice{},
        &InvoiceServiceLog{},
        &Part{},
        &ServiceLog{},
        &ServiceLogParts{},
    )
}
