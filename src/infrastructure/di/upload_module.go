package di

import (
	filesUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/sys/files"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	uploadController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/upload"
)

type UploadModule struct {
	Controller uploadController.IUploadController
	UseCase    files.ISysFilesRepository
	Repository files.ISysFilesRepository
}

func setupUploadModule(appContext *ApplicationContext) error {
	// Initialize repositories
	filesRepo := files.NewSysFilesRepository(appContext.DB, appContext.Logger)

	// Initialize use cases
	filesUC := filesUseCase.NewSysFilesUseCase(filesRepo, appContext.Logger)

	// Initialize controllers
	uploadController := uploadController.NewAuthController(filesUC, appContext.Logger)
	appContext.UploadModule = UploadModule{
		Controller: uploadController,
		UseCase:    filesUC,
		Repository: filesRepo,
	}
	return nil
}
