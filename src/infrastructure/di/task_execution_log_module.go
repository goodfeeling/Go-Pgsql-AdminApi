package di

import (
	taskExecutionLogUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	taskExecutionLogController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/task_execution_log"
	wsHandler "github.com/gbrayhan/microservices-go/src/infrastructure/ws/handler/task_execution_log"
)

type TaskExecutionLogModule struct {
	Controller taskExecutionLogController.ITaskExecutionLogController
	UseCase    taskExecutionLogUseCase.ITaskExecutionLogService
	Repository task_execution_log.ITaskExecutionLogRepository
	WsHandler  *wsHandler.LogHandler
}

func setupTaskExecutionLogModule(appContext *ApplicationContext) error {

	taskExecutionLogRepository := task_execution_log.NewTaskExecutionLogRepository(appContext.DB, appContext.Logger)
	// Initialize use cases
	services := taskExecutionLogUseCase.NewTaskExecutionLogUseCase(
		taskExecutionLogRepository,
		appContext.Logger)

	// Initialize websocket handler
	wsHandler := wsHandler.NewLogHandler(services, appContext.Logger, appContext.WsRouter)

	// Initialize controllers
	ctrl := taskExecutionLogController.NewTaskExecutionLogController(services, appContext.Logger)

	appContext.TaskExecutionLogModule = TaskExecutionLogModule{
		Controller: ctrl,
		UseCase:    services,
		Repository: taskExecutionLogRepository,
		WsHandler:  wsHandler,
	}
	return nil
}
