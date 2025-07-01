package role

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Role struct {
	ID            int64
	Name          string
	ParentID      int64
	DefaultRouter string
	Status        bool
	Order         int64
	Label         string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type SearchResultRole struct {
	Data       *[]Role `json:"data"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_page"`
}

type IRoleService interface {
	GetAll() (*[]Role, error)
	GetByID(id int) (*Role, error)
	Create(newRole *Role) (*Role, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*Role, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultRole, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*Role, error)
}
