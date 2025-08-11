package task_execution_log

import (
	"fmt"

	"github.com/gbrayhan/microservices-go/src/domain"
	taskExecutionLogkDomain "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	taskExecutionLogRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	"go.uber.org/zap"
)

type ITaskExecutionLogService interface {
	GetByID(id int) (*taskExecutionLogkDomain.TaskExecutionLog, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[taskExecutionLogkDomain.TaskExecutionLog], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}

type TaskExecutionLogUseCase struct {
	taskExecutionLogRepository taskExecutionLogRepo.ITaskExecutionLogRepository
	Logger                     *logger.Logger
}

func NewTaskExecutionLogUseCase(
	taskExecutionLogRepository taskExecutionLogRepo.ITaskExecutionLogRepository,
	loggerInstance *logger.Logger) ITaskExecutionLogService {
	return &TaskExecutionLogUseCase{
		taskExecutionLogRepository: taskExecutionLogRepository,
		Logger:                     loggerInstance,
	}
}

func (s *TaskExecutionLogUseCase) GetByID(id int) (*taskExecutionLogkDomain.TaskExecutionLog, error) {
	s.Logger.Info("Getting task by ID", zap.Int("id", id))
	return s.taskExecutionLogRepository.GetByID(id)
}

func (s *TaskExecutionLogUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting task", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.taskExecutionLogRepository.Delete(ids)
}

func (s *TaskExecutionLogUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[taskExecutionLogkDomain.TaskExecutionLog], error) {
	s.Logger.Info("Searching tasks with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.taskExecutionLogRepository.SearchPaginated(filters)
}

func (s *TaskExecutionLogUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching tasks by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.taskExecutionLogRepository.SearchByProperty(property, searchText)
}
