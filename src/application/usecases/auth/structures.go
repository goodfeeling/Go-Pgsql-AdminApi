package auth

import (
	"time"

	roleDomain "github.com/gbrayhan/microservices-go/src/domain/sys/role"
)

type LoginUser struct {
	UserName string
	Password string
}

type DataUserAuthenticated struct {
	UserName        string                 `json:"user_name"`
	NickName        string                 `json:"nick_name"`
	Email           string                 `json:"email"`
	Status          bool                   `json:"status"`
	UUID            string                 `json:"uuid"`
	ID              int64                  `json:"id"`
	Phone           string                 `json:"phone"`
	HeaderImg       string                 `json:"header_img"`
	CurrentRoleId   int64                  `json:"current_role_id"`
	CurrentRoleName string                 `json:"current_role_name"`
	Roles           []*roleDomain.RoleTree `json:"roles"`
}

type DataSecurityAuthenticated struct {
	JWTAccessToken            string    `json:"jwtAccessToken"`
	JWTRefreshToken           string    `json:"jwtRefreshToken"`
	ExpirationAccessDateTime  time.Time `json:"expirationAccessDateTime"`
	ExpirationRefreshDateTime time.Time `json:"expirationRefreshDateTime"`
}

type SecurityAuthenticatedUser struct {
	UserInfo DataUserAuthenticated     `json:"userinfo"`
	Security DataSecurityAuthenticated `json:"security"`
}

type SecurityRegisterUser struct {
	Data DataUserAuthenticated `json:"data"`
}

type RegisterUser struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
