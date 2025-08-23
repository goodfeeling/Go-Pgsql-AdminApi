package di

import (
	configUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/config"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/config"

	configController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/config"
)

type ConfigModule struct {
	Controller configController.IConfigController
	UseCase    configUseCase.ISysConfigService
	Repository config.ConfigRepositoryInterface
}

func setupConfigModule(appContext *ApplicationContext) error {
	// Initialize repositories
	configRepository := config.NewConfigRepository(appContext.DB, appContext.Logger)
	// Initialize use cases
	configUC := configUseCase.NewSysConfigUseCase(
		configRepository,
		appContext.Repositories.DictionaryRepository,
		appContext.Logger)
	// Initialize controllers
	configController := configController.NewConfigController(configUC, appContext.Logger)

	appContext.ConfigModule = ConfigModule{
		Controller: configController,
		UseCase:    configUC,
		Repository: configRepository,
	}
	return nil
}
