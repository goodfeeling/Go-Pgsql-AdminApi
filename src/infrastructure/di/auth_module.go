// #file:/root/myproject/microapp/src/infrastructure/di/modules/auth_module.go
package di

import (
	authUseCase "github.com/gbrayhan/microservices-go/src/application/services/auth"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
)

type AuthModule struct {
	Controller             authController.IAuthController
	UseCase                authUseCase.IAuthUseCase
	UserRepository         user.UserRepositoryInterface
	RoleRepository         role.ISysRolesRepository
	JwtBlacklistRepository jwt_blacklist.JwtBlacklistRepository
}

func setupAuthModule(appContext *ApplicationContext) error {
	// Initialize use cases
	authUC := authUseCase.NewAuthUseCase(
		appContext.Repositories.UserRepository,
		appContext.Repositories.RoleRepository,
		appContext.JWTService,
		appContext.Logger,
		appContext.Repositories.JwtBlacklistRepository,
		appContext.RedisClient)

	// Initialize controllers
	authController := authController.NewAuthController(authUC, appContext.Logger)

	appContext.AuthModule = AuthModule{
		Controller:             authController,
		UseCase:                authUC,
		UserRepository:         appContext.Repositories.UserRepository,
		JwtBlacklistRepository: appContext.Repositories.JwtBlacklistRepository,
		RoleRepository:         appContext.Repositories.RoleRepository,
	}
	return nil
}
