package menu

import (
	"strconv"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"

	"github.com/gbrayhan/microservices-go/src/domain"
	menuDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	"go.uber.org/zap"
)

type ISysMenuService interface {
	GetAll() ([]*menuDomain.MenuTree, error)
	GetByID(id int) (*menuDomain.Menu, error)
	Create(newMenu *menuDomain.Menu) (*menuDomain.Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*menuDomain.Menu, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuDomain.Menu], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuDomain.Menu, error)
	GetTreeMenus() (*menuDomain.MenuNode, error)
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

func (s *SysMenuUseCase) GetAll() ([]*menuDomain.MenuTree, error) {
	s.Logger.Info("Getting all menus")
	menus, err := s.sysMenuRepository.GetAll()
	if err != nil {
		return nil, err
	}
	menuMap := make(map[int]*menuDomain.MenuTree)
	var roots []*menuDomain.MenuTree

	// First traversal: Create all nodes and put them into the map.
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
			Level:       []int{},
			Children:    []*menuDomain.MenuTree{},
		}
		menuMap[item.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, item := range *menus {
		node := menuMap[item.ID]
		if item.ParentID == 0 {
			node.Level = []int{item.ID}
			roots = append(roots, node)
		} else {
			if parentNode, exists := menuMap[item.ParentID]; exists {
				// path handle
				node.Level = append(node.Level, parentNode.Level...)
				node.Level = append(node.Level, item.ID)

				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}
	return roots, nil
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

// GetTreeRoles implements ISysRoleService.
func (s *SysMenuUseCase) GetTreeMenus() (*menuDomain.MenuNode, error) {
	menus, err := s.sysMenuRepository.GetAll()
	if err != nil {
		return nil, err
	}
	menuMap := make(map[int]*menuDomain.MenuNode)
	var roots []*menuDomain.MenuNode

	// First traversal: Create all nodes and put them into the map.
	for _, item := range *menus {
		node := &menuDomain.MenuNode{
			ID:       strconv.Itoa(item.ID),
			Name:     item.Name,
			Key:      strconv.Itoa(item.ID),
			Path:     []int{},
			Children: []*menuDomain.MenuNode{},
		}
		menuMap[item.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, item := range *menus {
		node := menuMap[item.ID]
		if item.ParentID == 0 {
			node.Path = []int{item.ID}
			roots = append(roots, node)
		} else {
			if parentNode, exists := menuMap[item.ParentID]; exists {
				// path handle
				node.Path = append(node.Path, parentNode.Path...)
				node.Path = append(node.Path, item.ID)

				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}
	return &menuDomain.MenuNode{
		ID:       "0",
		Name:     "根节点",
		Key:      "0",
		Children: roots,
	}, nil
}
