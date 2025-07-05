package base_menu_btn

import (
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"gorm.io/gorm"
)

type SysBaseMenuBtn struct {
	ID            int64          `gorm:"column:id;primary_key" json:"id"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index:idx_sys_apis_deleted_at" json:"deletedAt,omitempty"`
	Name          string         `gorm:"column:name" json:"name,omitempty"`
	Desc          string         `gorm:"column:desc" json:"desc,omitempty"`
	SysBaseMenuID int64          `gorm:"column:sys_base_menu_id" json:"sysBaseMenuId,omitempty"`
}

func (SysBaseMenuBtn) TableName() string {
	return "sys_apis"
}

type ISysBaseMenuRepository interface {
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSysBaseMenuRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysBaseMenuRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
