package auth

import (
	"errors"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	jwtBlacklistDomain "github.com/gbrayhan/microservices-go/src/domain/jwt_blacklist"
	domainRole "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	sharedUtil "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IAuthUseCase interface {
	Login(username, password string) (*domainUser.User, *AuthTokens, *domainRole.Role, error)
	Logout(jwtToken string) (*domain.CommonResponse[string], error)
	Register(user RegisterUser) (*domain.CommonResponse[SecurityRegisterUser], error)
	AccessTokenByRefreshToken(refreshToken string) (*domainUser.User, *AuthTokens, error)
	SwitchRole(userId int, roleId int64) (*domainUser.User, *AuthTokens, *domainRole.Role, error)
}

type AuthUseCase struct {
	UserRepository         user.UserRepositoryInterface
	RoleRepository         role.ISysRolesRepository
	JWTService             security.IJWTService
	Logger                 *logger.Logger
	jwtBlacklistRepository jwtBlacklistDomain.IJwtBlacklistService
}

func NewAuthUseCase(
	userRepository user.UserRepositoryInterface,
	RoleRepository role.ISysRolesRepository,
	jwtService security.IJWTService,
	loggerInstance *logger.Logger,
	jwtBlacklistRepository jwtBlacklistDomain.IJwtBlacklistService,
) IAuthUseCase {
	return &AuthUseCase{
		UserRepository:         userRepository,
		RoleRepository:         RoleRepository,
		JWTService:             jwtService,
		Logger:                 loggerInstance,
		jwtBlacklistRepository: jwtBlacklistRepository,
	}
}

type AuthTokens struct {
	AccessToken               string
	RefreshToken              string
	ExpirationAccessDateTime  time.Time
	ExpirationRefreshDateTime time.Time
}

func (s *AuthUseCase) SwitchRole(userId int, roleId int64) (*domainUser.User, *AuthTokens, *domainRole.Role, error) {
	s.Logger.Info("User switch attempt", zap.Int("userId", userId))
	user, err := s.UserRepository.GetByID(int(userId))
	if err != nil {
		s.Logger.Error("Error getting user for switch", zap.Error(err), zap.Int("userId", userId))
		return nil, nil, nil, err
	}
	if user.ID == 0 {
		s.Logger.Warn("Login failed: user not found", zap.Int("userId", userId))
		return nil, nil, nil, domainErrors.NewAppError(errors.New("user don't no found"), domainErrors.NotAuthorized)
	}
	role, err := s.RoleRepository.GetByID(int(roleId))
	if err != nil {
		s.Logger.Error("Error getting role for switch", zap.Error(err), zap.Int("roleId", int(roleId)))
		return nil, nil, nil, err
	}
	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "access")
	if err != nil {
		s.Logger.Error("Error generating access token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}
	refreshTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "refresh")
	if err != nil {
		s.Logger.Error("Error generating refresh token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}

	s.Logger.Info("User login successful", zap.Int("userId", userId))
	return user, authTokens, role, nil
}

func (s *AuthUseCase) Login(username, password string) (*domainUser.User, *AuthTokens, *domainRole.Role, error) {
	s.Logger.Info("User login attempt", zap.String("username", username))
	user, err := s.UserRepository.GetByUsername(username)
	if err != nil {
		s.Logger.Error("Error getting user for login", zap.Error(err), zap.String("username", username))
		return nil, nil, nil, err
	}
	if user.ID == 0 {
		s.Logger.Warn("Login failed: user not found", zap.String("username", username))
		return nil, nil, nil, domainErrors.NewAppError(errors.New("username or password does not match"), domainErrors.NotAuthorized)
	}

	isAuthenticated := sharedUtil.CheckPasswordHash(password, user.HashPassword)
	if !isAuthenticated {
		s.Logger.Warn("Login failed: invalid password", zap.String("username", username))
		return nil, nil, nil, domainErrors.NewAppError(errors.New("username or password does not match"), domainErrors.NotAuthorized)
	}
	var role domainRole.Role
	var roleId int64
	if len(user.Roles) > 0 {
		roleId = user.Roles[0].ID
		role = user.Roles[0]
	}
	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "access")
	if err != nil {
		s.Logger.Error("Error generating access token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}
	refreshTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "refresh")
	if err != nil {
		s.Logger.Error("Error generating refresh token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}

	s.Logger.Info("User login successful", zap.String("username", username), zap.Int64("userID", user.ID))
	return user, authTokens, &role, nil
}

func (s *AuthUseCase) AccessTokenByRefreshToken(refreshToken string) (*domainUser.User, *AuthTokens, error) {
	s.Logger.Info("Refreshing access token")
	claimsMap, err := s.JWTService.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		s.Logger.Error("Error verifying refresh token", zap.Error(err))
		return nil, nil, err
	}
	userID := int(claimsMap["id"].(float64))
	user, err := s.UserRepository.GetByID(userID)
	if err != nil {
		s.Logger.Error("Error getting user for token refresh", zap.Error(err), zap.Int("userID", userID))
		return nil, nil, err
	}
	roleId := int64(claimsMap["role_id"].(float64))
	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "access")
	if err != nil {
		s.Logger.Error("Error generating new access token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, err
	}

	var expTime = int64(claimsMap["exp"].(float64))

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		RefreshToken:              refreshToken,
		ExpirationRefreshDateTime: time.Unix(expTime, 0),
	}

	s.Logger.Info("Access token refreshed successfully", zap.Int64("userID", user.ID))
	return user, authTokens, nil
}

// Register implements IAuthUseCase.
func (s *AuthUseCase) Register(user RegisterUser) (*domain.CommonResponse[SecurityRegisterUser], error) {
	// user is exist
	whereCondition := make(map[string]interface{}, 3)
	whereCondition["user_name"] = user.UserName
	dbUser, err := s.UserRepository.GetOneByMap(whereCondition)
	if err != nil {
		return nil, err
	}
	if dbUser.ID != 0 {
		return nil,
			domainErrors.NewAppError(errors.New("The user already exists"), domainErrors.UserExists)
	}
	userRepo := domainUser.User{
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
	}
	// password to has
	hash, err := sharedUtil.StringToHash(user.Password)
	if err != nil {
		return &domain.CommonResponse[SecurityRegisterUser]{}, err
	}
	userRepo.HashPassword = string(hash)

	// generate uuid
	userRepo.UUID = uuid.New().String()
	userRepo.Status = true

	res, err := s.UserRepository.Create(&userRepo)

	return &domain.CommonResponse[SecurityRegisterUser]{
		Data: SecurityRegisterUser{
			Data: DataUserAuthenticated{
				ID:       res.ID,
				UUID:     res.UUID,
				UserName: res.UserName,
				NickName: res.NickName,
				Email:    res.Email,
				Status:   res.Status,
			},
		},
		Status:  0,
		Message: "success",
	}, nil

}

func (s *AuthUseCase) Logout(jwtToken string) (*domain.CommonResponse[string], error) {
	var err error
	// check exist
	exist, err := s.jwtBlacklistRepository.IsJwtInBlacklist(jwtToken)
	if err != nil {
		return nil, domainErrors.NewAppError(err, domainErrors.TokenError)
	}
	if exist {
		return nil, domainErrors.NewAppError(errors.New("the user logout already"), domainErrors.TokenError)
	}
	err = s.jwtBlacklistRepository.AddToBlacklist(jwtToken)
	if err != nil {
		return nil, domainErrors.NewAppError(err, domainErrors.TokenError)
	}
	return &domain.CommonResponse[string]{Data: "true", Status: 0, Message: "success"}, nil
}
