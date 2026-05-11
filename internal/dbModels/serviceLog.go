package dbModels

import "time"

type ServiceLog struct {
    CustomerAmplifierId uint // TODO: FK -> CustomerAmplifiers
    TimeStart           time.Time
    TimeEnd             time.Time
    Description         string
    Notes               string
}
