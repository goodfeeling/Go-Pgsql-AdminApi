package auth

import "time"

type LoginUser struct {
	UserName string
	Password string
}

type DataUserAuthenticated struct {
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	Status    bool   `json:"status"`
	UUID      string `json:"uuid"`
	ID        int    `json:"id"`
	Phone     string `json:"phone"`
	HeaderImg string `json:"header_img"`
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
