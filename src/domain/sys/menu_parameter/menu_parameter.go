package menu_parameter

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type MenuParameter struct {
	ID        int       `json:"id"`
	MenuID    string    `json:"menu_id"`
	Type      string    `json:"type"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IMenuParameterService interface {
	GetAll() (*[]MenuParameter, error)
	GetByID(id int) (*MenuParameter, error)
	Create(newMenu *MenuParameter) (*MenuParameter, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*MenuParameter, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[MenuParameter], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*MenuParameter, error)
}
