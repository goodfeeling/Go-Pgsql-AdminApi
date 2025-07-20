package di

import (
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"gorm.io/gorm"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
)

type UserContext struct {
	UserController userController.IUserController
	UserUseCase    userUseCase.IUserUseCase
}

func setupUserDependencies(db *gorm.DB, logger *logger.Logger) (userController.IUserController, userUseCase.IUserUseCase, error) {
	userRepo := user.NewUserRepository(db, logger)
	userUC := userUseCase.NewUserUseCase(userRepo, logger)
	userController := userController.NewUserController(userUC, logger)

	return userController, userUC, nil
}
