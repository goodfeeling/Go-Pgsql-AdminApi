package upload

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
)

type IUploadUseCase interface {
}

type UploadUseCase struct {
	sysFilesRepository files.ISysFilesRepository
	Logger             *logger.Logger
}

func NewMedicineUseCase(sysFilesRepository files.ISysFilesRepository, loggerInstance *logger.Logger) IUploadUseCase {
	return &UploadUseCase{
		sysFilesRepository: sysFilesRepository,
		Logger:             loggerInstance,
	}
}
