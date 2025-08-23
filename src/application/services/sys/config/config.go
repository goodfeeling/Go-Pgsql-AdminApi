package config

import (
	configDomain "github.com/gbrayhan/microservices-go/src/domain/sys/config"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	configRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/config"
	dictionaryRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"

	"go.uber.org/zap"
)

type ISysConfigService interface {
	GetConfigByGroup() (*[]configDomain.GroupConfig, error)
	Update(module string, dataMap map[string]interface{}) error
	GetConfigByModule(module string) (*[]configDomain.Config, error)
}

type SysConfigUseCase struct {
	sysConfigRepository     configRepo.ConfigRepositoryInterface
	sysDictionaryRepository dictionaryRepo.DictionaryRepositoryInterface
	Logger                  *logger.Logger
}

func NewSysConfigUseCase(
	sysConfigRepository configRepo.ConfigRepositoryInterface,
	sysDictionaryRepository dictionaryRepo.DictionaryRepositoryInterface,
	loggerInstance *logger.Logger) ISysConfigService {
	return &SysConfigUseCase{
		sysConfigRepository:     sysConfigRepository,
		sysDictionaryRepository: sysDictionaryRepository,
		Logger:                  loggerInstance,
	}
}

// Update implements ISysConfigService.
func (s *SysConfigUseCase) Update(module string, userMap map[string]interface{}) error {
	s.Logger.Info("Updating config")

	for key, value := range userMap {
		configValue, ok := value.(string)
		if !ok {
			// 可选：记录日志或处理类型断言失败的情况
			continue
		}
		err := s.sysConfigRepository.UpdateByModule(module, key, configValue)
		if err != nil {
			// 可选：记录日志或处理更新失败的情况
			continue
		}
	}
	return nil
}

// GetConfigByGroup implements ISysConfigService.
func (s *SysConfigUseCase) GetConfigByGroup() (*[]configDomain.GroupConfig, error) {
	s.Logger.Info("Get config to group")
	list, err := s.sysConfigRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var groups []configDomain.GroupConfig
	for _, item := range *list {
		var group *configDomain.GroupConfig
		for i, g := range groups {
			if g.Name == item.Module {
				group = &groups[i]
				break
			}
		}
		if group == nil {
			group = &configDomain.GroupConfig{
				Name:    item.Module,
				Configs: []configDomain.Config{},
			}
			groups = append(groups, *group)
		}
	}

	for _, item := range *list {
		for i, g := range groups {
			if g.Name == item.Module {

				if item.ConfigType == "select" {
					dictData, err := s.sysDictionaryRepository.GetByType(item.ConfigKey)
					if err != nil {
						item.SelectOptions = nil
					} else {
						item.SelectOptions = dictData.Details
					}
				}
				groups[i].Configs = append(groups[i].Configs, item)

				break
			}
		}
	}

	return &groups, nil
}

// GetConfigByModule implements ISysConfigService.
func (s *SysConfigUseCase) GetConfigByModule(module string) (*[]configDomain.Config, error) {
	s.Logger.Info("get config by module", zap.String("module", module))
	return s.sysConfigRepository.GetConfigByModule(module)
}
