package dbModels

import "gorm.io/gorm"

type InvoiceServiceLog struct {
    gorm.Model
    InvoiceId    uint // TODO: FK -> Invoices
    ServiceLogId uint // TODO: FK -> ServiceLogs
}
