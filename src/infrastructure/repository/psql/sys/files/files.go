package files

import (
	"fmt"
	"os"
	"time"

	filesDomain "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysFiles struct {
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`

	FileName string `gorm:"column:file_name;size:191;" json:"fileName"`
	FileMD5  string `gorm:"column:file_md5;size:191;" json:"fileMD5"`
	FilePath string `gorm:"column:file_path;size:191;" json:"filePath"`
	ID       int64  `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
}

func (SysFiles) TableName() string {
	return "sys_files"
}

type ISysFilesRepository interface {
	Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

// Create implements ISysFilesRepository.
func (r *Repository) Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error) {
	r.Logger.Info("Upload new file", zap.String("filename", data.FileName))
	fileRepository := fromDomainMapper(data)
	txDb := r.DB.Create(fileRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating user", zap.Error(err), zap.String("filename", data.FileName))
	}
	r.Logger.Info("Successfully add file", zap.String("filename", data.FileName), zap.Int("id", int(fileRepository.ID)))
	return fileRepository.toDomainMapper(), err
}

func NewSysFilesRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysFilesRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}

func fromDomainMapper(u *filesDomain.SysFiles) *SysFiles {
	return &SysFiles{
		FileName: u.FileName,
		FileMD5:  u.FileMD5,
		FilePath: u.FilePath,
	}
}

func (u *SysFiles) toDomainMapper() *filesDomain.SysFiles {
	return &filesDomain.SysFiles{
		ID:       u.ID,
		FileName: u.FileName,
		FileMD5:  u.FileMD5,
		FileUrl:  fmt.Sprintf("%s/%s", os.Getenv("APP_URL"), u.FilePath),
	}
}
