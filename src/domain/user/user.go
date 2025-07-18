package user

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	roleDomain "github.com/gbrayhan/microservices-go/src/domain/sys/role"
)

type User struct {
	ID            int64
	UUID          string
	UserName      string
	NickName      string
	Email         string
	Status        bool
	HashPassword  string
	HeaderImg     string
	Phone         string
	OriginSetting string
	Password      string
	RoleId        int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Roles         []roleDomain.Role
}
type SearchResultUser struct {
	Data       *[]User `json:"data"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_page"`
}

type IUserService interface {
	GetAll() (*[]User, error)
	GetByID(id int) (*User, error)
	Create(newUser *User) (*User, error)
	Delete(id int) error
	Update(id int64, userMap map[string]interface{}) (*User, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*User, error)
}
