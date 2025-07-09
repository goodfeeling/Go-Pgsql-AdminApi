package dictionary_detail

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type Dictionary struct {
	ID              int       `json:"id"`
	Label           string    `json:"label"`
	Value           string    `json:"value"`
	Extend          string    `json:"extend"`
	Status          int16     `json:"status"`
	Sort            int8      `json:"sort"`
	SysDictionaryID int64     `json:"sys_dictionary_id"`
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
