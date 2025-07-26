package di

import (
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"

	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
)

type UserModule struct {
	Controller userController.IUserController
	UseCase    userUseCase.IUserUseCase
	Repository user.UserRepositoryInterface
}

func setupUserModule(appContext *ApplicationContext) error {
	// Initialize use cases
	userUC := userUseCase.NewUserUseCase(
		appContext.Repositories.UserRepository,
		appContext.Repositories.UserRoleRepository, appContext.Logger)

	// Initialize controllers
	userController := userController.NewUserController(userUC, appContext.Logger)

	appContext.UserModule = UserModule{
		Controller: userController,
		UseCase:    userUC,
		Repository: appContext.Repositories.UserRepository,
	}
	return nil
}
