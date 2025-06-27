package authorities

import (
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"gorm.io/gorm"
)

type SysAuthority struct {
	AuthorityID     int64      `gorm:"primaryKey;column:authority_id;type:numeric(20,0)"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;index"`
	AuthorityName   string     `gorm:"column:authority_name"`
	ParentID        *int64     `gorm:"column:parent_id;type:numeric(20,0)"`
	DefaultRouter   string     `gorm:"column:default_router"`
	AuthorityStatus bool       `gorm:"column:authority_status"`
	AuthorityOrder  *int64     `gorm:"column:authority_order;type:numeric(10,0)"`
	Label           string     `gorm:"column:label"`
	Description     string     `gorm:"column:description"`
}

func (SysAuthority) TableName() string {
	return "sys_apis"
}

type ISysAuthorityRepository interface {
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSysAuthorityRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysAuthorityRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
