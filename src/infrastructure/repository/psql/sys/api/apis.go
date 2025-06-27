package apis

import (
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"gorm.io/gorm"
)

type SysApi struct {
	ID          int64      `gorm:"primaryKey;column:id;type:numeric(20,0)"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`
	Path        string     `gorm:"column:path"`
	Description string     `gorm:"column:description"`
	ApiGroup    string     `gorm:"column:api_group"`
	Method      string     `gorm:"column:method"`
}

func (SysApi) TableName() string {
	return "sys_apis"
}

type ISysFilesRepository interface {
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSysFilesRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysFilesRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
