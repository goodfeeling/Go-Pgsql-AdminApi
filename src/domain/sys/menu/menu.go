package menu

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Menu struct {
	ID          int       `json:"id"`
	Path        string    `json:"path"`
	MenuGroup   string    `json:"api_group"`
	Method      string    `json:"method"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IMenuService interface {
	GetAll() (*[]Menu, error)
	GetByID(id int) (*Menu, error)
	Create(newMenu *Menu) (*Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*Menu, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[Menu], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*Menu, error)
}
