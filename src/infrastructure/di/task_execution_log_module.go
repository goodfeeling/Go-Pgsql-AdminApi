package di

import (
	taskExecutionLogUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	taskExecutionLogController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/task_execution_log"
)

type TaskExecutionLogModule struct {
	Controller taskExecutionLogController.ITaskExecutionLogController
	UseCase    taskExecutionLogUseCase.ITaskExecutionLogService
	Repository task_execution_log.ITaskExecutionLogRepository
}

func setupTaskExecutionLogModule(appContext *ApplicationContext) error {
	// Initialize repositories
	taskExecutionLogRepository := task_execution_log.NewTaskExecutionLogRepository(appContext.DB, appContext.Logger)
	// Initialize use cases
	services := taskExecutionLogUseCase.NewTaskExecutionLogUseCase(
		taskExecutionLogRepository,
		appContext.Logger)
	// Initialize controllers
	taskExecutionLogController := taskExecutionLogController.NewITaskExecutionLogController(services, appContext.Logger)

	appContext.TaskExecutionLogModule = TaskExecutionLogModule{
		Controller: taskExecutionLogController,
		UseCase:    services,
		Repository: taskExecutionLogRepository,
	}
	return nil
}
