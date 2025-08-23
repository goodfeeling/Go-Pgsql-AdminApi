package config

import (
	configDomain "github.com/gbrayhan/microservices-go/src/domain/sys/config"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	configRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/config"
	"go.uber.org/zap"
)

type ISysConfigService interface {
	GetConfigByGroup() (*[]configDomain.GroupConfig, error)
	Update(dataMap map[string]interface{}) (*configDomain.Config, error)
	GetConfigByModule(module string) (*[]configDomain.Config, error)
}

type SysConfigUseCase struct {
	sysConfigRepository configRepo.ConfigRepositoryInterface
	Logger              *logger.Logger
}

func NewSysConfigUseCase(
	sysConfigRepository configRepo.ConfigRepositoryInterface,
	loggerInstance *logger.Logger) ISysConfigService {
	return &SysConfigUseCase{
		sysConfigRepository: sysConfigRepository,
		Logger:              loggerInstance,
	}
}

// Update implements ISysConfigService.
func (s *SysConfigUseCase) Update(userMap map[string]interface{}) (*configDomain.Config, error) {
	s.Logger.Info("Updating config")

	for _, value := range userMap {
		config, ok := value.(configDomain.Config)
		if !ok {
			// 可选：记录日志或处理类型断言失败的情况
			continue
		}
		s.sysConfigRepository.Update(&config)
	}
	return nil, nil
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
				groups[i].Configs = append(groups[i].Configs, item) // 修改这里，通过索引操作
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
