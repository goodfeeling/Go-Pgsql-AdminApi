package files

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type SysFiles struct {
	ID             int64     `json:"id"`
	FileName       string    `json:"file_name"`
	FileMD5        string    `json:"file_md5"`
	FilePath       string    `json:"file_path"`
	FileUrl        string    `json:"file_url"`
	StorageEngine  string    `json:"storage_engine"`
	FileOriginName string    `json:"file_origin_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ISysFilesService interface {
	Create(data *SysFiles) (*SysFiles, error)
	GetAll() (*[]SysFiles, error)
	GetByID(id int) (*SysFiles, error)
	Delete(ids []int64) error
	Update(id int, userMap map[string]interface{}) (*SysFiles, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[SysFiles], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*SysFiles, error)
}
