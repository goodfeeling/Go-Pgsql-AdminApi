package api

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Api struct {
	ID          int       `json:"id"`
	Path        string    `json:"path"`
	ApiGroup    string    `json:"api_group"`
	Method      string    `json:"method"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IApiService interface {
	GetAll() (*[]Api, error)
	GetByID(id int) (*Api, error)
	Create(newApi *Api) (*Api, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*Api, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[Api], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*Api, error)
	GetApisGroup() (*[]GroupApiItem, error)
}
type GroupApiItem struct {
	GroupName string          `json:"title"`
	GroupKey  string          `json:"key"`
	Children  []*GroupApiItem `json:"children"`
}
