package di

import (
	scheduledTaskUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/scheduled_task"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/scheduled_task"
	scheduledTaskController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/scheduled_task"
)

type ScheduledTaskModule struct {
	Controller scheduledTaskController.IScheduledTaskDetailController
	UseCase    scheduledTaskUseCase.IScheduledTaskService
	Repository scheduled_task.IScheduledTaskRepository
}

func setupScheduledTaskModule(appContext *ApplicationContext) error {
	// Initialize use cases
	apiUC := scheduledTaskUseCase.NewScheduledTaskUseCase(
		appContext.Repositories.ScheduledTaskRepository,
		appContext.Logger)
	// Initialize controllers
	scheduledTaskController := scheduledTaskController.NewScheduledTaskDetailController(apiUC, appContext.Logger)

	appContext.ScheduledTaskModule = ScheduledTaskModule{
		Controller: scheduledTaskController,
		UseCase:    apiUC,
		Repository: appContext.Repositories.ScheduledTaskRepository,
	}
	return nil
}
