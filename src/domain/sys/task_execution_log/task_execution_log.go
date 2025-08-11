package task_execution_log

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type TaskExecutionLog struct {
	ID              int       `json:"id"`
	TaskID          uint      `json:"task_id"`
	ExecuteTime     time.Time `json:"execute_time"`
	ExecuteResult   int       `json:"execute_result"`   // 1-成功, 0-失败
	ExecuteDuration *int      `json:"execute_duration"` // 执行耗时(毫秒)
	ErrorMessage    *string   `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ITaskExecutionLogService interface {
	GetByID(id int) (*TaskExecutionLog, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[TaskExecutionLog], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}
