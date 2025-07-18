package menu

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Menu struct {
	ID          int         `json:"id"`
	MenuLevel   int         `json:"menu_level"`
	ParentID    int         `json:"parent_id"`
	Path        string      `json:"path"`
	Name        string      `json:"name"`
	Hidden      int16       `json:"hidden"`
	Component   string      `json:"component"`
	Sort        int8        `json:"sort"`
	ActiveName  string      `json:"active_name"`
	KeepAlive   int16       `json:"keep_alive"`
	DefaultMenu int16       `json:"default_menu"`
	Title       string      `json:"title"`
	Icon        string      `json:"icon"`
	CloseTab    int16       `json:"close_tab"`
	MenuGroupId int         `json:"menu_group_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Level       []int       `json:"level"`
	Children    []*MenuTree `json:"children"`
}

type MenuNode struct {
	ID       string      `json:"value"`
	Name     string      `json:"title"`
	Key      string      `json:"key"`
	Path     []int       `json:"path"`
	Children []*MenuNode `json:"children"`
}

type MenuTree struct {
	ID          int               `json:"id"`
	MenuLevel   int               `json:"menu_level"`
	ParentID    int               `json:"parent_id"`
	Path        string            `json:"path"`
	Name        string            `json:"name"`
	Hidden      int16             `json:"hidden"`
	Component   string            `json:"component"`
	Sort        int8              `json:"sort"`
	ActiveName  string            `json:"active_name"`
	KeepAlive   int16             `json:"keep_alive"`
	DefaultMenu int16             `json:"default_menu"`
	Title       string            `json:"title"`
	Icon        string            `json:"icon"`
	CloseTab    int16             `json:"close_tab"`
	CreatedAt   domain.CustomTime `json:"created_at"`
	UpdatedAt   domain.CustomTime `json:"updated_at"`
	Level       []int             `json:"level"`
	Children    []*MenuTree       `json:"children"`
}

type MenuGroup struct {
	Name  string      `json:"name"`
	Path  string      `json:"path"`
	Items []*MenuTree `json:"items"`
}

type IMenuService interface {
	GetAll(groupId int) ([]*MenuTree, error)
	GetByID(id int) (*Menu, error)
	Create(newMenu *Menu) (*Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*Menu, error)
	GetOneByMap(userMap map[string]interface{}) (*Menu, error)
	GetUserMenus(userId int) ([]*MenuGroup, error)
}
