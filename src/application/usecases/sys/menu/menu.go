package menu

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	menuDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	menuGroupRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	roleMenuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	userRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"go.uber.org/zap"
)

type ISysMenuService interface {
	GetAll(groupId int) ([]*menuDomain.MenuTree, error)
	GetByID(id int) (*menuDomain.Menu, error)
	Create(newMenu *menuDomain.Menu) (*menuDomain.Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*menuDomain.Menu, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuDomain.Menu], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuDomain.Menu, error)
	GetUserMenus(userId int) ([]*menuDomain.MenuGroup, error)
}

type SysMenuUseCase struct {
	sysMenuRepository      menuRepo.MenuRepositoryInterface
	userRepository         userRepo.UserRepositoryInterface
	sysRoleMenuRepository  roleMenuRepo.ISysRoleMenuRepository
	sysMenuGroupRepository menuGroupRepo.MenuGroupRepositoryInterface
	Logger                 *logger.Logger
}

func NewSysMenuUseCase(
	sysMenuRepository menuRepo.MenuRepositoryInterface,
	sysRoleMenuRepository roleMenuRepo.ISysRoleMenuRepository,
	userRepository userRepo.UserRepositoryInterface,
	sysMenuGroupRepository menuGroupRepo.MenuGroupRepositoryInterface,
	loggerInstance *logger.Logger,
) ISysMenuService {
	return &SysMenuUseCase{
		sysMenuRepository:      sysMenuRepository,
		userRepository:         userRepository,
		sysRoleMenuRepository:  sysRoleMenuRepository,
		sysMenuGroupRepository: sysMenuGroupRepository,
		Logger:                 loggerInstance,
	}
}

func (s *SysMenuUseCase) GetAll(groupId int) ([]*menuDomain.MenuTree, error) {
	s.Logger.Info("Getting all menus")
	menus, err := s.sysMenuRepository.GetAll(groupId)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, ""), nil

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

// GetUserMenus
func (s *SysMenuUseCase) GetUserMenus(userId int) ([]*menuDomain.MenuGroup, error) {
	s.Logger.Info("Getting user menus", zap.Int("userId", userId))
	userInfo, err := s.userRepository.GetByID(userId)
	if err != nil {
		return nil, err
	}
	roleId := userInfo.RoleId
	s.Logger.Info("Getting user menus", zap.Int64("roleId", roleId))
	roleMenus, err := s.sysRoleMenuRepository.GetByRoleId(roleId)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("Getting user menus", zap.Int("menusCount", len(roleMenus)))
	menus, err := s.sysMenuRepository.GetByIDs(roleMenus)
	if err != nil {
		return nil, err
	}
	groups, err := s.sysMenuGroupRepository.GetAll()
	if err != nil {
		return nil, err
	}
	menuGroupMap := make(map[int][]menuDomain.Menu)
	for _, menu := range *menus {
		menuGroupMap[menu.MenuGroupId] = append(menuGroupMap[menu.MenuGroupId], menuDomain.Menu{
			ID:          menu.ID,
			Name:        menu.Name,
			Title:       menu.Title,
			Path:        menu.Path,
			Component:   menu.Component,
			Icon:        menu.Icon,
			CloseTab:    menu.CloseTab,
			DefaultMenu: menu.DefaultMenu,
			Hidden:      menu.Hidden,
			Sort:        menu.Sort,
			ActiveName:  menu.ActiveName,
			KeepAlive:   menu.KeepAlive,
			MenuLevel:   menu.MenuLevel,
			ParentID:    menu.ParentID,
			CreatedAt:   menu.CreatedAt,
			UpdatedAt:   menu.UpdatedAt,
		})
	}
	menuGroup := make([]*menuDomain.MenuGroup, 0)
	for _, group := range *groups {
		menusForGroup := menuGroupMap[group.ID]
		menuGroup = append(menuGroup, &menuDomain.MenuGroup{
			Name:  group.Name,
			Path:  group.Path,
			Items: buildMenuTree(&menusForGroup, group.Path),
		})
	}
	return menuGroup, nil
}

// buildMenuTree
func buildMenuTree(menus *[]menuDomain.Menu, groupPath string) []*menuDomain.MenuTree {

	menuMap := make(map[int]*menuDomain.MenuTree)
	var roots []*menuDomain.MenuTree

	//  traversal: Establish parent-child relationships.
	for _, item := range *menus {
		node := &menuDomain.MenuTree{
			ID:          item.ID,
			Path:        item.Path,
			Name:        item.Name,
			ParentID:    item.ParentID,
			DefaultMenu: item.DefaultMenu,
			Hidden:      item.Hidden,
			MenuLevel:   item.MenuLevel,
			CloseTab:    item.CloseTab,
			KeepAlive:   item.KeepAlive,
			Icon:        item.Icon,
			Title:       item.Title,
			Sort:        item.Sort,
			ActiveName:  item.ActiveName,
			Component:   item.Component,
			CreatedAt:   domain.CustomTime{Time: item.CreatedAt},
			UpdatedAt:   domain.CustomTime{Time: item.UpdatedAt},
			Level:       []int{item.ID},
			Children:    []*menuDomain.MenuTree{},
		}
		menuMap[item.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, item := range *menus {
		node := menuMap[item.ID]
		if item.ParentID == 0 {
			roots = append(roots, node)
		} else if parentNode, exists := menuMap[item.ParentID]; exists {
			// path handle
			node.Level = append(node.Level, parentNode.Level...)

			// api get user menu handle
			if groupPath != "" {
				node.Path = parentNode.Path + "/" + node.Path
			}

			parentNode.Children = append(parentNode.Children, node)
		} else {
			// 父节点不存在，作为孤儿节点加入根节点列表
			node.Level = []int{item.ID}
			roots = append(roots, node)
		}
	}

	if groupPath != "" {
		// 只在最终叶子节点添加 groupPath 前缀，并避免重复拼接
		for _, node := range menuMap {
			if len(node.Children) == 0 {
				// 只添加 groupPath 前缀，不重复拼接
				node.Path = "/" + groupPath + "/" + node.Path
			}
		}
	}

	return roots
}
