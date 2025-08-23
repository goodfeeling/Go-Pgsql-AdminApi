package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainConfig "github.com/gbrayhan/microservices-go/src/domain/sys/config"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SysConfig 系统配置表
type SysConfig struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ConfigKey   string    `gorm:"column:config_key;type:varchar(100);not null;uniqueIndex:idx_sys_config_key" json:"config_key" binding:"required"`
	ConfigValue string    `gorm:"column:config_value;type:text" json:"config_value"`
	ConfigType  string    `gorm:"column:config_type;type:varchar(20);default:string;check:config_type IN ('string', 'number', 'boolean', 'json', 'array')" json:"config_type"`
	Module      string    `gorm:"column:module;type:varchar(50);index:idx_sys_config_module" json:"module"`
	EnvType     string    `gorm:"column:env_type;type:varchar(20);default:default;index:idx_sys_config_env" json:"env_type"`
	Sort        int       `gorm:"column:sort;type:int;default:0" json:"sort"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (SysConfig) TableName() string {
	return "sys_config"
}

var ColumnsConfigMapping = map[string]string{
	"id":          "id",
	"path":        "path",
	"configName":  "config_name",
	"description": "description",
	"configGroup": "config_group",
	"method":      "method",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

// ConfigRepositoryInterface defines the interface for config repository operations
type ConfigRepositoryInterface interface {
	GetAll() (*[]domainConfig.Config, error)
	Create(configDomain *domainConfig.Config) (*domainConfig.Config, error)
	GetByID(id int) (*domainConfig.Config, error)
	Update(configDomain *domainConfig.Config) (*domainConfig.Config, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainConfig.Config], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(configMap map[string]interface{}) (*domainConfig.Config, error)
	GetConfigByModule(module string) (*[]domainConfig.Config, error)
	UpdateByModule(module, configKey, configValue string) error
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewConfigRepository(db *gorm.DB, loggerInstance *logger.Logger) ConfigRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainConfig.Config, error) {
	var configs []SysConfig
	if err := r.DB.Order("sort ASC").Find(&configs).Error; err != nil {
		r.Logger.Error("Error getting all configs", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all configs", zap.Int("count", len(configs)))
	return arrayToDomainMapper(&configs), nil
}

func (r *Repository) Create(configDomain *domainConfig.Config) (*domainConfig.Config, error) {
	r.Logger.Info("Creating new config", zap.String("ConfigKey", configDomain.ConfigKey))
	configRepository := fromDomainMapper(configDomain)
	txDb := r.DB.Create(configRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating config", zap.Error(err), zap.String("ConfigKey", configDomain.ConfigKey))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainConfig.Config{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainConfig.Config{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created config", zap.String("ConfigKey", configDomain.ConfigKey), zap.Int("id", int(configRepository.ID)))
	return configRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainConfig.Config, error) {
	var config SysConfig
	err := r.DB.Where("id = ?", id).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Config not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting config by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainConfig.Config{}, err
	}
	r.Logger.Info("Successfully retrieved config by ID", zap.Int("id", id))
	return config.toDomainMapper(), nil
}

func (r *Repository) Update(configDomain *domainConfig.Config) (*domainConfig.Config, error) {
	var configObj SysConfig
	configObj.ID = configDomain.ID
	err := r.DB.Model(&configObj).Updates(fromDomainMapper(configDomain)).Error
	if err != nil {
		r.Logger.Error("Error updating config", zap.Error(err), zap.Int64("id", configDomain.ID))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainConfig.Config{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainConfig.Config{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainConfig.Config{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", configDomain.ID).First(&configObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated config", zap.Error(err), zap.Int64("id", configDomain.ID))
		return &domainConfig.Config{}, err
	}
	r.Logger.Info("Successfully updated config", zap.Int64("id", configDomain.ID))
	return configObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&SysConfig{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting config", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Config not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted config", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainConfig.Config], error) {
	query := r.DB.Model(&SysConfig{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsConfigMapping[field]
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
			column := ColumnsConfigMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsConfigMapping[dateFilter.Field]
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
			column := ColumnsConfigMapping[sortField]
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

	var configs []SysConfig
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&configs).Error; err != nil {
		r.Logger.Error("Error searching configs", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainConfig.Config]{
		Data:       arrayToDomainMapper(&configs),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched configs",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsConfigMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysConfig{}).
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

func (u *SysConfig) toDomainMapper() *domainConfig.Config {
	return &domainConfig.Config{
		ID:          u.ID,
		ConfigKey:   u.ConfigKey,
		ConfigType:  u.ConfigType,
		ConfigValue: u.ConfigValue,
		Sort:        u.Sort,
		Module:      u.Module,
		EnvType:     u.EnvType,
		CreatedAt:   domain.CustomTime{Time: u.CreatedAt},
		UpdatedAt:   domain.CustomTime{Time: u.UpdatedAt},
	}
}

func fromDomainMapper(u *domainConfig.Config) *SysConfig {
	return &SysConfig{
		ID:          u.ID,
		ConfigKey:   u.ConfigKey,
		ConfigType:  u.ConfigType,
		ConfigValue: u.ConfigValue,
		Sort:        u.Sort,
		Module:      u.Module,
		EnvType:     u.EnvType,
	}
}

func arrayToDomainMapper(configs *[]SysConfig) *[]domainConfig.Config {
	configsDomain := make([]domainConfig.Config, len(*configs))
	for i, config := range *configs {
		configsDomain[i] = *config.toDomainMapper()
	}
	return &configsDomain
}

func (r *Repository) GetOneByMap(configMap map[string]interface{}) (*domainConfig.Config, error) {
	var configRepository SysConfig
	tx := r.DB.Limit(1)
	for key, value := range configMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&configRepository).Error; err != nil {
		return &domainConfig.Config{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return configRepository.toDomainMapper(), nil
}

func (r *Repository) GetConfigByModule(module string) (*[]domainConfig.Config, error) {
	var configs []SysConfig
	if err := r.DB.Where("module = ?", module).Find(&configs).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMapper(&configs), nil
}

func (r *Repository) UpdateByModule(module, configKey, configValue string) error {
	envType := os.Getenv("ENV_TYPE")
	err := r.DB.
		Model(&SysConfig{}).
		Where("module = ? and config_key = ? and env_type = ?", module, configKey, envType).
		Update("config_value", configValue).Error
	if err != nil {
		r.Logger.Error("Error updating config", zap.Error(err), zap.String("configKey", configKey))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)

	}
	fmt.Println(module, configKey, envType)
	r.Logger.Info("Successfully updated config")
	return nil
}
