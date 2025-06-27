package files

import (
	filesDomain "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	"go.uber.org/zap"
)

type ISysFilesService interface {
	Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error)
}

type SysFilesUseCase struct {
	sysFilesRepository files.ISysFilesRepository
	Logger             *logger.Logger
}

// Create implements ISysFilesService.
func (s *SysFilesUseCase) Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error) {
	s.Logger.Info("Getting file by filename", zap.String("filename", data.FileName))
	return s.sysFilesRepository.Create(data)
}

func NewSysFilesUseCase(sysFilesRepository files.ISysFilesRepository, loggerInstance *logger.Logger) ISysFilesService {
	return &SysFilesUseCase{
		sysFilesRepository: sysFilesRepository,
		Logger:             loggerInstance,
	}
}
