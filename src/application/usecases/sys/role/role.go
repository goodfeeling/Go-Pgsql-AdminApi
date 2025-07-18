package role

import (
	"strconv"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	casbinRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/casbin_rule"
	roleRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	roleMenuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"

	"github.com/gbrayhan/microservices-go/src/domain"
	roleDomain "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	"go.uber.org/zap"
)

type ISysRoleService interface {
	GetAll() ([]*roleDomain.RoleTree, error)
	GetByID(id int) (*roleDomain.Role, error)
	GetByName(name string) (*roleDomain.Role, error)
	Create(newRole *roleDomain.Role) (*roleDomain.Role, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*roleDomain.Role, error)
	SearchPaginated(filters domain.DataFilters) (*roleDomain.SearchResultRole, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*roleDomain.Role, error)
	GetTreeRoles() (*roleDomain.RoleNode, error)

	GetRoleMenuIds(id int64) ([]int, error)
	UpdateRoleMenuIds(id int, updateMap map[string]any) error

	GetApiRuleList(roleId int) ([]string, error)
	BindApiRule(roleId int, updateMap map[string]interface{}) error
}

type SysRoleUseCase struct {
	sysRoleRepository     roleRepo.ISysRolesRepository
	sysRoleMenuRepository roleMenuRepo.ISysRoleMenuRepository
	casbinRuleRepo        casbinRepo.ICasbinRuleRepository

	Logger *logger.Logger
}

func NewSysRoleUseCase(sysRoleRepository roleRepo.ISysRolesRepository, sysRoleMenuRepository roleMenuRepo.ISysRoleMenuRepository, casbinRuleRepo casbinRepo.ICasbinRuleRepository, loggerInstance *logger.Logger) ISysRoleService {
	return &SysRoleUseCase{
		sysRoleRepository:     sysRoleRepository,
		sysRoleMenuRepository: sysRoleMenuRepository,
		casbinRuleRepo:        casbinRuleRepo,
		Logger:                loggerInstance,
	}
}

func (s *SysRoleUseCase) GetAll() ([]*roleDomain.RoleTree, error) {

	s.Logger.Info("Getting all roles")
	roles, err := s.sysRoleRepository.GetAll()
	if err != nil {
		return nil, err
	}
	roleMap := make(map[int64]*roleDomain.RoleTree)
	var roots []*roleDomain.RoleTree

	// First traversal: Create all nodes and put them into the map.
	for _, role := range *roles {
		node := &roleDomain.RoleTree{
			ID:            role.ID,
			Name:          role.Name,
			ParentID:      role.ParentID,
			DefaultRouter: role.DefaultRouter,
			Status:        role.Status,
			Order:         role.Order,
			Label:         role.Label,
			Description:   role.Description,
			CreatedAt:     role.CreatedAt,
			UpdatedAt:     role.UpdatedAt,
			Path:          []int64{},
			Children:      []*roleDomain.RoleTree{},
		}
		roleMap[role.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, role := range *roles {
		node := roleMap[role.ID]
		if role.ParentID == 0 {
			node.Path = []int64{role.ID}
			roots = append(roots, node)
		} else {
			if parentNode, exists := roleMap[role.ParentID]; exists {
				// path handle
				node.Path = append(node.Path, parentNode.Path...)
				node.Path = append(node.Path, role.ID)

				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}
	return roots, nil
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

// GetTreeRoles implements ISysRoleService.
func (s *SysRoleUseCase) GetTreeRoles() (*roleDomain.RoleNode, error) {
	roles, err := s.sysRoleRepository.GetAll()
	if err != nil {
		return nil, err
	}
	roleMap := make(map[int64]*roleDomain.RoleNode)
	var roots []*roleDomain.RoleNode

	// First traversal: Create all nodes and put them into the map.
	for _, role := range *roles {
		id := strconv.Itoa(int(role.ID))
		node := &roleDomain.RoleNode{
			ID:       id,
			Name:     role.Name,
			Key:      id,
			Path:     []int64{},
			Children: []*roleDomain.RoleNode{},
		}
		roleMap[role.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, role := range *roles {
		node := roleMap[role.ID]
		if role.ParentID == 0 {
			node.Path = []int64{role.ID}
			roots = append(roots, node)
		} else {
			if parentNode, exists := roleMap[role.ParentID]; exists {
				// path handle
				node.Path = append(node.Path, parentNode.Path...)
				node.Path = append(node.Path, role.ID)
				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}
	return &roleDomain.RoleNode{
		ID:       "0",
		Name:     "根节点",
		Key:      "0",
		Children: roots,
	}, nil
}

func (s *SysRoleUseCase) GetRoleMenuIds(id int64) ([]int, error) {
	return s.sysRoleMenuRepository.GetByRoleId(id)
}

func (s *SysRoleUseCase) UpdateRoleMenuIds(id int, updateMap map[string]any) error {
	return s.sysRoleMenuRepository.Insert(id, updateMap)
}
func (s *SysRoleUseCase) GetApiRuleList(roleId int) ([]string, error) {
	return s.casbinRuleRepo.GetByRoleId(roleId)

}
func (s *SysRoleUseCase) BindApiRule(roleId int, updateMap map[string]interface{}) error {
	return s.casbinRuleRepo.Insert(roleId, updateMap)
}
