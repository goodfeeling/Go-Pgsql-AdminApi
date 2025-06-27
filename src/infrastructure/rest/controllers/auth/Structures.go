package auth

import "time"

type LoginRequest struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type SecurityData struct {
	JWTAccessToken            string    `json:"jwtAccessToken"`
	JWTRefreshToken           string    `json:"jwtRefreshToken"`
	ExpirationAccessDateTime  time.Time `json:"expirationAccessDateTime"`
	ExpirationRefreshDateTime time.Time `json:"expirationRefreshDateTime"`
}

type RegisterRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
