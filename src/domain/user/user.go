package user

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type User struct {
	ID            int
	UUID          string
	UserName      string
	NickName      string
	Email         string
	Status        bool
	HashPassword  string
	HeaderImg     string
	AuthorityId   int64
	Phone         string
	OriginSetting string
	Password      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type SearchResultUser struct {
	Data       *[]User
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type IUserService interface {
	GetAll() (*[]User, error)
	GetByID(id int) (*User, error)
	Create(newUser *User) (*User, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*User, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*User, error)
}
