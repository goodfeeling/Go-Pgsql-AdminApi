package files

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	filesDomain "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysFiles struct {
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`

	FileName       string `gorm:"column:file_name;size:191;" json:"fileName"`
	FileMD5        string `gorm:"column:file_md5;size:191;" json:"fileMD5"`
	FilePath       string `gorm:"column:file_path;size:191;" json:"filePath"`
	StorageEngine  string `gorm:"column:storage_engine;size:10;" json:"storageEngine"`
	FileOriginName string `gorm:"column:file_origin_name;size:191;" json:"fileOriginName"`
	ID             int64  `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
}

func (SysFiles) TableName() string {
	return "sys_files"
}

var ColumnsSysFilesMapping = map[string]string{
	"id":             "id",
	"path":           "path",
	"dictionaryName": "dictionary_name",
	"selectedDictId": "dictionary_group",
	"method":         "method",
	"createdAt":      "created_at",
	"updatedAt":      "updated_at",
}

type ISysFilesRepository interface {
	Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error)
	GetAll() (*[]filesDomain.SysFiles, error)
	GetByID(id int) (*filesDomain.SysFiles, error)
	Update(id int, dictionaryMap map[string]interface{}) (*filesDomain.SysFiles, error)
	Delete(ids []int64) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[filesDomain.SysFiles], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(dictionaryMap map[string]interface{}) (*filesDomain.SysFiles, error)
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
		FileName:       u.FileName,
		FileMD5:        u.FileMD5,
		FilePath:       u.FilePath,
		StorageEngine:  u.StorageEngine,
		FileOriginName: u.FileOriginName,
	}
}

func (u *SysFiles) toDomainMapper() *filesDomain.SysFiles {
	return &filesDomain.SysFiles{
		ID:             u.ID,
		FileName:       u.FileName,
		FileMD5:        u.FileMD5,
		FileUrl:        fmt.Sprintf("%s/%s", os.Getenv("APP_URL"), u.FilePath),
		StorageEngine:  u.StorageEngine,
		FileOriginName: u.FileOriginName,
		CreatedAt:      *u.CreatedAt,
		UpdatedAt:      *u.UpdatedAt,
	}
}

func (r *Repository) GetAll() (*[]filesDomain.SysFiles, error) {
	var files []SysFiles
	if err := r.DB.Find(&files).Error; err != nil {
		r.Logger.Error("Error getting all files", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all files", zap.Int("count", len(files)))
	return arrayToDomainMapper(&files), nil
}

func (r *Repository) GetByID(id int) (*filesDomain.SysFiles, error) {
	var dictionary SysFiles
	err := r.DB.Where("id = ?", id).First(&dictionary).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("SysFiles not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting dictionary by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &filesDomain.SysFiles{}, err
	}
	r.Logger.Info("Successfully retrieved dictionary by ID", zap.Int("id", id))
	return dictionary.toDomainMapper(), nil
}

func (r *Repository) Update(id int, dictionaryMap map[string]interface{}) (*filesDomain.SysFiles, error) {
	var dictionaryObj SysFiles
	dictionaryObj.ID = int64(id)
	err := r.DB.Model(&dictionaryObj).
		Updates(dictionaryMap).Error
	if err != nil {
		r.Logger.Error("Error updating dictionary", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &filesDomain.SysFiles{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &filesDomain.SysFiles{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &filesDomain.SysFiles{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&dictionaryObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated dictionary", zap.Error(err), zap.Int("id", id))
		return &filesDomain.SysFiles{}, err
	}
	r.Logger.Info("Successfully updated dictionary", zap.Int("id", id))
	return dictionaryObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int64) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&SysFiles{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting file", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Api not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted file", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[filesDomain.SysFiles], error) {
	query := r.DB.Model(&SysFiles{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsSysFilesMapping[field]
					if column != "" {
						query = query.Where(column+" ILIKE ?", "%"+value+"%")
					}
				}
			}
		}
	}

	// Apply exact matches
	for field, values := range filters.Matches {
		if len(values) > 0 {
			column := ColumnsSysFilesMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsSysFilesMapping[dateFilter.Field]
		if column != "" {
			if dateFilter.Start != nil {
				query = query.Where(column+" >= ?", dateFilter.Start)
			}
			if dateFilter.End != nil {
				query = query.Where(column+" <= ?", dateFilter.End)
			}
		}
	}

	// Apply sorting
	if len(filters.SortBy) > 0 && filters.SortDirection.IsValid() {
		for _, sortField := range filters.SortBy {
			column := ColumnsSysFilesMapping[sortField]
			if column != "" {
				query = query.Order(column + " " + string(filters.SortDirection))
			}
		}
	}

	// Count total records
	var total int64
	clonedQuery := query
	clonedQuery.Count(&total)

	// Apply pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	offset := (filters.Page - 1) * filters.PageSize

	var files []SysFiles
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&files).Error; err != nil {
		r.Logger.Error("Error searching files", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[filesDomain.SysFiles]{
		Data:       arrayToDomainMapper(&files),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched files",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsSysFilesMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysFiles{}).
		Distinct(column).
		Where(column+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(column, &coincidences).Error; err != nil {
		r.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	r.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(coincidences)))

	return &coincidences, nil
}

func arrayToDomainMapper(files *[]SysFiles) *[]filesDomain.SysFiles {
	filesDomain := make([]filesDomain.SysFiles, len(*files))
	for i, dictionary := range *files {
		filesDomain[i] = *dictionary.toDomainMapper()
	}
	return &filesDomain
}

func (r *Repository) GetOneByMap(dictionaryMap map[string]interface{}) (*filesDomain.SysFiles, error) {
	var dictionaryRepository SysFiles
	tx := r.DB.Limit(1)
	for key, value := range dictionaryMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&dictionaryRepository).Error; err != nil {
		return &filesDomain.SysFiles{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return dictionaryRepository.toDomainMapper(), nil
}
