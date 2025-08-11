package di

import (
	eventHandler "github.com/gbrayhan/microservices-go/src/application/event/handler"
	eventModel "github.com/gbrayhan/microservices-go/src/application/event/model"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/services/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"

	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
)

type UserModule struct {
	Controller userController.IUserController
	UseCase    userUseCase.IUserUseCase
	Repository user.UserRepositoryInterface
}

func setupUserModule(appContext *ApplicationContext) error {

	// Initialize event
	appContext.EventBus.Subscribe(
		eventModel.UserRegisteredEventType, eventHandler.NewNotificationEventHandler())

	// Initialize use cases
	userUC := userUseCase.NewUserUseCase(
		appContext.Repositories.UserRepository,
		appContext.Repositories.UserRoleRepository,
		appContext.EventBus,
		appContext.Logger)

	// Initialize controllers
	userController := userController.NewUserController(userUC, appContext.Logger)

	appContext.UserModule = UserModule{
		Controller: userController,
		UseCase:    userUC,
		Repository: appContext.Repositories.UserRepository,
	}
	return nil
}
