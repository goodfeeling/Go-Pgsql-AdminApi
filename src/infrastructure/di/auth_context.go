package di

import (
	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"gorm.io/gorm"
)

type AuthContext struct {
	AuthController authController.IAuthController
	AuthUseCase    authUseCase.IAuthUseCase
	UserRepository user.UserRepositoryInterface
}

func setupAuthDependencies(db *gorm.DB, logger *logger.Logger) (authController.IAuthController, authUseCase.IAuthUseCase, error) {
	jwtService := security.NewJWTService()
	userRepo := user.NewUserRepository(db, logger)
	jwtBlackListRepo := jwt_blacklist.NewUJwtBlacklistRepository(db)

	authUC := authUseCase.NewAuthUseCase(userRepo, jwtService, logger, jwtBlackListRepo)
	authController := authController.NewAuthController(authUC, logger)

	return authController, authUC, nil
}
