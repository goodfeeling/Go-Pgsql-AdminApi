package api

import (
	"fmt"
	"strings"

	"github.com/gbrayhan/microservices-go/src/domain"
	apiDomain "github.com/gbrayhan/microservices-go/src/domain/sys/api"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	apiRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"
	dictionaryRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ISysApiService interface {
	GetAll() (*[]apiDomain.Api, error)
	GetByID(id int) (*apiDomain.Api, error)
	Create(newApi *apiDomain.Api) (*apiDomain.Api, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*apiDomain.Api, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[apiDomain.Api], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*apiDomain.Api, error)
	GetApisGroup() (*[]apiDomain.GroupApiItem, error)
	SynchronizeRouterToApi(router gin.RoutesInfo) (*int, error)
}

type SysApiUseCase struct {
	sysApiRepository     apiRepo.ApiRepositoryInterface
	dictionaryRepository dictionaryRepo.DictionaryRepositoryInterface
	Logger               *logger.Logger
}

func NewSysApiUseCase(
	sysApiRepository apiRepo.ApiRepositoryInterface,
	dictionaryRepository dictionaryRepo.DictionaryRepositoryInterface,
	loggerInstance *logger.Logger) ISysApiService {
	return &SysApiUseCase{
		sysApiRepository:     sysApiRepository,
		dictionaryRepository: dictionaryRepository,
		Logger:               loggerInstance,
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

func (s *SysApiUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting api", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.sysApiRepository.Delete(ids)
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

// Get one api by map
func (s *SysApiUseCase) GetOneByMap(userMap map[string]interface{}) (*apiDomain.Api, error) {
	return s.sysApiRepository.GetOneByMap(userMap)
}

// GetApisGroup
func (s *SysApiUseCase) GetApisGroup() (*[]apiDomain.GroupApiItem, error) {
	apis, err := s.sysApiRepository.GetAll()
	if err != nil {
		return nil, err
	}

	dictionary, err := s.dictionaryRepository.GetByType("api_group")
	if err != nil {
		return nil, err
	}

	groups := make([]apiDomain.GroupApiItem, len(*dictionary.Details))
	for i, item := range *dictionary.Details {
		groupApis := make([]*apiDomain.GroupApiItem, 0)
		for _, api := range *apis {
			if api.ApiGroup == item.Value {
				groupApis = append(groupApis, &apiDomain.GroupApiItem{
					GroupKey:  fmt.Sprintf("%v---%v", api.Path, api.Method),
					GroupName: api.Description,
					Children:  nil,
				})
			}
		}
		groups[i] = apiDomain.GroupApiItem{
			GroupName:       item.Label,
			GroupKey:        fmt.Sprintf("0---%v", i),
			DisableCheckbox: len(groupApis) == 0,
			Children:        groupApis,
		}
	}
	return &groups, nil
}

func (c *SysApiUseCase) SynchronizeRouterToApi(routes gin.RoutesInfo) (*int, error) {
	count := 0
	for _, route := range routes {
		if c.shouldSyncRoute(route.Path) {
			apiModel := &apiRepo.SysApi{
				Path:        route.Path,
				Method:      route.Method,
				Description: c.generateDescription(route.Path, route.Method),
				ApiGroup:    "other",
			}

			ok, err := c.sysApiRepository.Upsert(apiModel)
			if err != nil {
				c.Logger.Error("Failed to sync route",
					zap.String("path", route.Path),
					zap.String("method", route.Method),
					zap.Error(err))
				continue
			}
			if ok {
				count++
			}
		}
	}
	return &count, nil
}

func (a *SysApiUseCase) shouldSyncRoute(path string) bool {
	// 排除一些系统路由
	excludePaths := []string{"/swagger", "/health"}
	for _, exclude := range excludePaths {
		if strings.HasPrefix(path, exclude) {
			return false
		}
	}
	return true
}

func (a *SysApiUseCase) generateDescription(path, method string) string {
	// 根据路径和方法生成描述
	switch method {
	case "GET":
		if strings.Contains(path, "/:id") {
			return "获取单个资源"
		}
		return "获取资源列表"
	case "POST":
		return "创建资源"
	case "PUT":
		return "更新资源"
	case "DELETE":
		return "删除资源"
	default:
		return "API接口"
	}
}
