package dbModels

type ServiceLogParts struct {
    ServiceLogId uint // TODO: FK -> ServiceLogs
    PartId       uint // TODO: FK -> Parts
    Qty          uint
}
