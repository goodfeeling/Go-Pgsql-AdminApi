package operation_records

import "time"

type SysOperationRecord struct {
	ID           int64
	IP           string
	Method       string
	Path         string
	Status       int8
	Latency      int64
	Agent        string
	ErrorMessage string
	Body         string
	Resp         string
	UserID       int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
type ISysOperationRecordService interface {
	Create(record *SysOperationRecord) error
}
