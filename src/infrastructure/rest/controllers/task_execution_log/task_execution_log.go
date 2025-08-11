package task_execution_log

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainTaskExecutionLog "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	taskExecutionLogRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type DeleteBatchTaskExecutionLogRequest struct {
	IDS []int `json:"ids"`
}

type ResponseTaskExecutionLog struct {
	ID              int               `json:"id"`
	TaskID          uint              `json:"task_id"`
	ExecuteTime     time.Time         `json:"execute_time"`
	ExecuteResult   int               `json:"execute_result"`
	ExecuteDuration *int              `json:"execute_duration"`
	ErrorMessage    *string           `json:"error_message"`
	CreatedAt       domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt       domain.CustomTime `json:"updated_at,omitempty"`
}
type ITaskExecutionLogController interface {
	GetTaskExecutionLogByID(ctx *gin.Context)
	DeleteTaskExecutionLog(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
	DeleteTaskExecutionLogs(ctx *gin.Context)
}
type TaskExecutionLogController struct {
	taskExecutionLogService domainTaskExecutionLog.ITaskExecutionLogService
	Logger                  *logger.Logger
}

func NewITaskExecutionLogController(taskExecutionLogService domainTaskExecutionLog.ITaskExecutionLogService, loggerInstance *logger.Logger) ITaskExecutionLogController {
	return &TaskExecutionLogController{taskExecutionLogService: taskExecutionLogService, Logger: loggerInstance}
}

// GetTaskExecutionLogByID
// @Summary get scheduled_task
// @Description get scheduled_task by id
// @Tags scheduled_task
// @Accept json
// @Produce json
// @Success 200 {object} ResponseTaskExecutionLog
// @Router /v1/taskExecutionLog/{id} [get]
func (c *TaskExecutionLogController) GetTaskExecutionLogByID(ctx *gin.Context) {
	taskExecutionLogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid taskExecutionLog ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("taskExecutionLog id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting taskExecutionLog by ID", zap.Int("id", taskExecutionLogID))
	taskExecutionLog, err := c.taskExecutionLogService.GetByID(taskExecutionLogID)
	if err != nil {
		c.Logger.Error("Error getting taskExecutionLog by ID", zap.Error(err), zap.Int("id", taskExecutionLogID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved taskExecutionLog by ID", zap.Int("id", taskExecutionLogID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(taskExecutionLog))
}

// DeleteTaskExecutionLog
// @Summary delete taskExecutionLog
// @Description delete taskExecutionLog by id
// @Tags taskExecutionLog
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/taskExecutionLog/{id} [delete]
func (c *TaskExecutionLogController) DeleteTaskExecutionLog(ctx *gin.Context) {
	taskExecutionLogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid taskExecutionLog ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting taskExecutionLog", zap.Int("id", taskExecutionLogID))
	err = c.taskExecutionLogService.Delete([]int{taskExecutionLogID})
	if err != nil {
		c.Logger.Error("Error deleting taskExecutionLog", zap.Error(err), zap.Int("id", taskExecutionLogID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("TaskExecutionLog deleted successfully", zap.Int("id", taskExecutionLogID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    taskExecutionLogID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// SearchTaskExecutionLogPageList
// @Summary search scheduled_task
// @Description search scheduled_task by query
// @Tags search scheduled_task
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseTaskExecutionLog]
// @Router /v1/taskExecutionLog/search [get]
func (c *TaskExecutionLogController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching scheduled_task with pagination")

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Build filters
	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse like filters
	likeFilters := make(map[string][]string)
	for field := range taskExecutionLogRepo.ColumnsTaskExecutionLogMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range taskExecutionLogRepo.ColumnsTaskExecutionLogMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	fmt.Println(matches)
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range taskExecutionLogRepo.ColumnsTaskExecutionLogMapping {
		startStr := ctx.Query(field + "_start")
		endStr := ctx.Query(field + "_end")

		if startStr != "" || endStr != "" {
			dateRange := domain.DateRangeFilter{Field: field}

			if startStr != "" {
				if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
					dateRange.Start = &startTime
				}
			}

			if endStr != "" {
				if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
					dateRange.End = &endTime
				}
			}

			dateRanges = append(dateRanges, dateRange)
		}
	}
	filters.DateRangeFilters = dateRanges

	// Parse sorting
	sortBy := ctx.QueryArray("sortBy")
	if len(sortBy) > 0 {
		filters.SortBy = sortBy
	} else {
		filters.SortBy = []string{}
	}

	sortDirection := domain.SortDirection(ctx.DefaultQuery("sortDirection", "asc"))
	if sortDirection.IsValid() {
		filters.SortDirection = sortDirection
	}

	result, err := c.taskExecutionLogService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching scheduled_task", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseTaskExecutionLog]
	response := controllers.NewCommonResponseBuilder[PageResult]().
		Data(PageResult{
			List:       arrayDomainToResponseMapper(result.Data),
			Total:      result.Total,
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalPages: result.TotalPages,
			Filters:    filters,
		}).
		Message("success").
		Status(0).
		Build()

	c.Logger.Info("Successfully searched scheduled_task",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

// SearchByProperty
// @Summary  search by property
// @Description search by property
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {array} []string
// @Router /v1/taskExecutionLog/search-property [get]
func (c *TaskExecutionLogController) SearchByProperty(ctx *gin.Context) {
	property := ctx.Query("property")
	searchText := ctx.Query("searchText")

	if property == "" || searchText == "" {
		c.Logger.Error("Missing property or searchText parameter")
		appError := domainErrors.NewAppError(errors.New("missing property or searchText parameter"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	// Validate property
	allowed := map[string]bool{
		"taskExecutionLogName": true,
		"email":                true,
		"firstName":            true,
		"lastName":             true,
		"status":               true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.taskExecutionLogService.SearchByProperty(property, searchText)
	if err != nil {
		c.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(*coincidences)))
	ctx.JSON(http.StatusOK, coincidences)
}

// DeleteTaskExecutionLogs
// @Summary delete logs
// @Description delete logs by id
// @Tags batch delete
// @Accept json
// @Produce json
// @Param book body DeleteBatchOperationRequest true  "JSON Data"
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/operation/delete-batch [post]
func (c *TaskExecutionLogController) DeleteTaskExecutionLogs(ctx *gin.Context) {
	c.Logger.Info("Creating new taskExecutionLog")
	var request DeleteBatchTaskExecutionLogRequest
	var err error
	if err = controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new taskExecutionLog", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting operation", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	err = c.taskExecutionLogService.Delete(request.IDS)
	if err != nil {
		c.Logger.Error("Error deleting operation", zap.Error(err), zap.String("ids", fmt.Sprintf("%v", request.IDS)))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Operation deleted successfully", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[[]int]{
		Data:    request.IDS,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// Mappers
func domainToResponseMapper(domainTaskExecutionLog *domainTaskExecutionLog.TaskExecutionLog) *ResponseTaskExecutionLog {

	return &ResponseTaskExecutionLog{
		ID:        domainTaskExecutionLog.ID,
		CreatedAt: domain.CustomTime{Time: domainTaskExecutionLog.CreatedAt},
		UpdatedAt: domain.CustomTime{Time: domainTaskExecutionLog.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(scheduled_task *[]domainTaskExecutionLog.TaskExecutionLog) *[]*ResponseTaskExecutionLog {
	res := make([]*ResponseTaskExecutionLog, len(*scheduled_task))
	for i, u := range *scheduled_task {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}
