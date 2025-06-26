package files

import (
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"gorm.io/gorm"
)

type SysFiles struct {
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`

	FileName string `gorm:"column:file_name;size:191;" json:"fileName"`
	FileMD5  string `gorm:"column:file_md5;size:191;" json:"fileMD5"`
	FilePath string `gorm:"column:file_path;size:191;" json:"filePath"`
	ID       int64  `gorm:"primaryKey;autoIncrement:false" json:"id"`
}

func (SysFiles) TableName() string {
	return "sys_files"
}

type ISysFilesRepository interface{}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSysFilesRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysFilesRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
