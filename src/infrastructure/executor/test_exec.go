package executor

import (
	"fmt"

	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
)

func CleanOldData(*domainScheduledTask.ScheduledTask) error {
	fmt.Println("say hello")
	return nil
}
