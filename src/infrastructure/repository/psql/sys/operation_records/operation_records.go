package operation_records

import (
	"time"

	operationRecordsDomain "github.com/gbrayhan/microservices-go/src/domain/sys/operation_records"
	"gorm.io/gorm"
)

// SysOperationRecord represents the sys_operation_records table structure.
type SysOperationRecord struct {
	ID           int64      `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	IP           string     `gorm:"column:ip" json:"ip,omitempty"`
	Method       string     `gorm:"column:method" json:"method,omitempty"`
	Path         string     `gorm:"column:path" json:"path,omitempty"`
	Status       int8       `gorm:"column:status" json:"status,omitempty"`
	Latency      int64      `gorm:"column:latency" json:"latency,omitempty"`
	Agent        string     `gorm:"column:agent" json:"agent,omitempty"`
	ErrorMessage string     `gorm:"column:error_message" json:"errorMessage,omitempty"`
	Body         string     `gorm:"column:body" json:"body,omitempty"`
	Resp         string     `gorm:"column:resp" json:"resp,omitempty"`
	UserID       int64      `gorm:"column:user_id" json:"userId,omitempty"`
}

func (*SysOperationRecord) TableName() string {
	return "sys_operation_records"
}

type ISysOperationRecordsRepository interface {
	Create(record *operationRecordsDomain.SysOperationRecord) error
}

type Repository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) ISysOperationRecordsRepository {
	return &Repository{DB: db}
}

// Create implements ISysOperationRecordsRepository.
func (r *Repository) Create(record *operationRecordsDomain.SysOperationRecord) error {
	recordRep := fromDomainMapper(record)
	txDB := r.DB.Create(recordRep)
	txErr := txDB.Error
	if txErr != nil {
		return txErr
	}
	return nil
}

func fromDomainMapper(u *operationRecordsDomain.SysOperationRecord) *SysOperationRecord {
	return &SysOperationRecord{
		CreatedAt:    &u.CreatedAt,
		IP:           u.IP,
		Method:       u.Method,
		Path:         u.Path,
		Status:       u.Status,
		Latency:      u.Latency,
		Agent:        u.Agent,
		ErrorMessage: u.ErrorMessage,
		Body:         u.Body,
		Resp:         u.Resp,
		UserID:       u.UserID,
	}
}
