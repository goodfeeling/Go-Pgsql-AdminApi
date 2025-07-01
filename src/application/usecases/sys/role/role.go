package role

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	roleRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"

	"github.com/gbrayhan/microservices-go/src/domain"
	roleDomain "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	"go.uber.org/zap"
)

type ISysRoleService interface {
	GetAll() (*[]roleDomain.Role, error)
	GetByID(id int) (*roleDomain.Role, error)
	GetByName(name string) (*roleDomain.Role, error)
	Create(newRole *roleDomain.Role) (*roleDomain.Role, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*roleDomain.Role, error)
	SearchPaginated(filters domain.DataFilters) (*roleDomain.SearchResultRole, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*roleDomain.Role, error)
}

type SysRoleUseCase struct {
	sysRoleRepository roleRepo.ISysRolesRepository
	Logger            *logger.Logger
}

func NewSysFilesUseCase(sysRoleRepository roleRepo.ISysRolesRepository, loggerInstance *logger.Logger) ISysRoleService {
	return &SysRoleUseCase{
		sysRoleRepository: sysRoleRepository,
		Logger:            loggerInstance,
	}
}

func (s *SysRoleUseCase) GetAll() (*[]roleDomain.Role, error) {
	s.Logger.Info("Getting all roles")
	return s.sysRoleRepository.GetAll()
}

func (s *SysRoleUseCase) GetByID(id int) (*roleDomain.Role, error) {
	s.Logger.Info("Getting role by ID", zap.Int("id", id))
	return s.sysRoleRepository.GetByID(id)
}

func (s *SysRoleUseCase) GetByName(name string) (*roleDomain.Role, error) {
	s.Logger.Info("Getting role by name", zap.String("name", name))
	return s.sysRoleRepository.GetByName(name)
}

func (s *SysRoleUseCase) Create(newRole *roleDomain.Role) (*roleDomain.Role, error) {
	s.Logger.Info("Creating new role", zap.String("name", newRole.Name))
	return s.sysRoleRepository.Create(newRole)
}

func (s *SysRoleUseCase) Delete(id int) error {
	s.Logger.Info("Deleting role", zap.Int("id", id))
	return s.sysRoleRepository.Delete(id)
}

func (s *SysRoleUseCase) Update(id int, userMap map[string]interface{}) (*roleDomain.Role, error) {
	s.Logger.Info("Updating role", zap.Int("id", id))
	return s.sysRoleRepository.Update(id, userMap)
}

func (s *SysRoleUseCase) SearchPaginated(filters domain.DataFilters) (*roleDomain.SearchResultRole, error) {
	s.Logger.Info("Searching roles with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysRoleRepository.SearchPaginated(filters)
}

func (s *SysRoleUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching role by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysRoleRepository.SearchByProperty(property, searchText)
}

func (s *SysRoleUseCase) GetOneByMap(userMap map[string]interface{}) (*roleDomain.Role, error) {
	return s.sysRoleRepository.GetOneByMap(userMap)
}
