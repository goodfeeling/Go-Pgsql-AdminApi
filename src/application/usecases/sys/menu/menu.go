package menu

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"

	"github.com/gbrayhan/microservices-go/src/domain"
	menuDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	"go.uber.org/zap"
)

type ISysMenuService interface {
	GetAll() (*[]menuDomain.Menu, error)
	GetByID(id int) (*menuDomain.Menu, error)
	Create(newMenu *menuDomain.Menu) (*menuDomain.Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*menuDomain.Menu, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuDomain.Menu], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuDomain.Menu, error)
}

type SysMenuUseCase struct {
	sysMenuRepository menuRepo.MenuRepositoryInterface
	Logger            *logger.Logger
}

func NewSysMenuUseCase(sysMenuRepository menuRepo.MenuRepositoryInterface, loggerInstance *logger.Logger) ISysMenuService {
	return &SysMenuUseCase{
		sysMenuRepository: sysMenuRepository,
		Logger:            loggerInstance,
	}
}

func (s *SysMenuUseCase) GetAll() (*[]menuDomain.Menu, error) {
	s.Logger.Info("Getting all roles")
	return s.sysMenuRepository.GetAll()
}

func (s *SysMenuUseCase) GetByID(id int) (*menuDomain.Menu, error) {
	s.Logger.Info("Getting menu by ID", zap.Int("id", id))
	return s.sysMenuRepository.GetByID(id)
}

func (s *SysMenuUseCase) Create(newMenu *menuDomain.Menu) (*menuDomain.Menu, error) {
	s.Logger.Info("Creating new menu", zap.String("path", newMenu.Path))
	return s.sysMenuRepository.Create(newMenu)
}

func (s *SysMenuUseCase) Delete(id int) error {
	s.Logger.Info("Deleting menu", zap.Int("id", id))
	return s.sysMenuRepository.Delete(id)
}

func (s *SysMenuUseCase) Update(id int, userMap map[string]interface{}) (*menuDomain.Menu, error) {
	s.Logger.Info("Updating menu", zap.Int("id", id))
	return s.sysMenuRepository.Update(id, userMap)
}

func (s *SysMenuUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuDomain.Menu], error) {
	s.Logger.Info("Searching menus with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysMenuRepository.SearchPaginated(filters)
}

func (s *SysMenuUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching menu by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysMenuRepository.SearchByProperty(property, searchText)
}

func (s *SysMenuUseCase) GetOneByMap(userMap map[string]interface{}) (*menuDomain.Menu, error) {
	return s.sysMenuRepository.GetOneByMap(userMap)
}
