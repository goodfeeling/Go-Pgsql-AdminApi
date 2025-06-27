package base_menus

import (
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"gorm.io/gorm"
)

type SysBaseMenu struct {
	ID          int64      `gorm:"primaryKey;column:id;type:numeric(20,0)"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`
	MenuLevel   *int64     `gorm:"column:menu_level;type:numeric(20,0)"`
	ParentID    *int64     `gorm:"column:parent_id;type:numeric(20,0)"`
	Path        string     `gorm:"column:path"`
	Name        string     `gorm:"column:name"`
	Hidden      int16      `gorm:"column:hidden"`
	Component   string     `gorm:"column:component"`
	Sort        int8       `gorm:"column:sort"`
	ActiveName  string     `gorm:"column:active_name"`
	KeepAlive   int16      `gorm:"column:keep_alive"`
	DefaultMenu int16      `gorm:"column:default_menu"`
	Title       string     `gorm:"column:title"`
	Icon        string     `gorm:"column:icon"`
	CloseTab    int16      `gorm:"column:close_tab"`
}

func (SysBaseMenu) TableName() string {
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
