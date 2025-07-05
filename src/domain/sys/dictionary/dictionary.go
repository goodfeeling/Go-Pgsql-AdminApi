package dictionary

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Dictionary struct {
	ID              int       `json:"id"`
	Path            string    `json:"path"`
	DictionaryGroup string    `json:"api_group"`
	Method          string    `json:"method"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type IDictionaryService interface {
	GetAll() (*[]Dictionary, error)
	GetByID(id int) (*Dictionary, error)
	Create(newDictionary *Dictionary) (*Dictionary, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*Dictionary, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[Dictionary], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*Dictionary, error)
}
