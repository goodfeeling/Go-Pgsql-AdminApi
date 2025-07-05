package api

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	apiRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"

	"github.com/gbrayhan/microservices-go/src/domain"
	apiDomain "github.com/gbrayhan/microservices-go/src/domain/sys/api"
	"go.uber.org/zap"
)

type ISysApiService interface {
	GetAll() (*[]apiDomain.Api, error)
	GetByID(id int) (*apiDomain.Api, error)
	Create(newApi *apiDomain.Api) (*apiDomain.Api, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*apiDomain.Api, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[apiDomain.Api], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*apiDomain.Api, error)
}

type SysApiUseCase struct {
	sysApiRepository apiRepo.ApiRepositoryInterface
	Logger           *logger.Logger
}

func NewSysApiUseCase(sysApiRepository apiRepo.ApiRepositoryInterface, loggerInstance *logger.Logger) ISysApiService {
	return &SysApiUseCase{
		sysApiRepository: sysApiRepository,
		Logger:           loggerInstance,
	}
}

func (s *SysApiUseCase) GetAll() (*[]apiDomain.Api, error) {
	s.Logger.Info("Getting all roles")
	return s.sysApiRepository.GetAll()
}

func (s *SysApiUseCase) GetByID(id int) (*apiDomain.Api, error) {
	s.Logger.Info("Getting api by ID", zap.Int("id", id))
	return s.sysApiRepository.GetByID(id)
}

func (s *SysApiUseCase) Create(newApi *apiDomain.Api) (*apiDomain.Api, error) {
	s.Logger.Info("Creating new api", zap.String("path", newApi.Path))
	return s.sysApiRepository.Create(newApi)
}

func (s *SysApiUseCase) Delete(id int) error {
	s.Logger.Info("Deleting api", zap.Int("id", id))
	return s.sysApiRepository.Delete(id)
}

func (s *SysApiUseCase) Update(id int, userMap map[string]interface{}) (*apiDomain.Api, error) {
	s.Logger.Info("Updating api", zap.Int("id", id))
	return s.sysApiRepository.Update(id, userMap)
}

func (s *SysApiUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[apiDomain.Api], error) {
	s.Logger.Info("Searching apis with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysApiRepository.SearchPaginated(filters)
}

func (s *SysApiUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching api by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysApiRepository.SearchByProperty(property, searchText)
}

func (s *SysApiUseCase) GetOneByMap(userMap map[string]interface{}) (*apiDomain.Api, error) {
	return s.sysApiRepository.GetOneByMap(userMap)
}
