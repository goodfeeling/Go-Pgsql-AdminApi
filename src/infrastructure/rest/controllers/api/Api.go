package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainApi "github.com/gbrayhan/microservices-go/src/domain/sys/api"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	apiRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewApiRequest struct {
	ID          int    `json:"id"`
	Path        string `json:"path" binding:"required"`
	ApiGroup    string `json:"api_group" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ResponseApi struct {
	ID          int               `json:"id"`
	Path        string            `json:"path"`
	ApiGroup    string            `json:"api_group"`
	Method      string            `json:"method"`
	Description string            `json:"description"`
	CreatedAt   domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt   domain.CustomTime `json:"updated_at,omitempty"`
}
type IApiController interface {
	NewApi(ctx *gin.Context)
	GetAllApis(ctx *gin.Context)
	GetApisByID(ctx *gin.Context)
	UpdateApi(ctx *gin.Context)
	DeleteApi(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
}
type ApiController struct {
	apiService domainApi.IApiService
	Logger     *logger.Logger
}

func NewApiController(apiService domainApi.IApiService, loggerInstance *logger.Logger) IApiController {
	return &ApiController{apiService: apiService, Logger: loggerInstance}
}

// CreateApi
// @Summary create api
// @Description create api
// @Tags api create
// @Accept json
// @Produce json
// @Param book body NewApiRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/api [post]
func (c *ApiController) NewApi(ctx *gin.Context) {
	c.Logger.Info("Creating new api")
	var request NewApiRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new api", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	apiModel, err := c.apiService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating api", zap.Error(err), zap.String("path", request.Path))
		_ = ctx.Error(err)
		return
	}
	apiResponse := controllers.NewCommonResponseBuilder[*ResponseApi]().
		Data(domainToResponseMapper(apiModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Api created successfully", zap.String("path", request.Path), zap.Int("id", int(apiModel.ID)))
	ctx.JSON(http.StatusOK, apiResponse)
}

// GetAllApis
// @Summary get all apis by
// @Description get  all apis by where
// @Tags apis
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainApi.Api]
// @Router /v1/api [get]
func (c *ApiController) GetAllApis(ctx *gin.Context) {
	c.Logger.Info("Getting all apis")
	apis, err := c.apiService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all apis", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all apis", zap.Int("count", len(*apis)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainApi.Api]{
		Data: apis,
	})
}

// GetApisByID
// @Summary get apis
// @Description get apis by id
// @Tags apis
// @Accept json
// @Produce json
// @Success 200 {object} ResponseApi
// @Router /v1/api/{id} [get]
func (c *ApiController) GetApisByID(ctx *gin.Context) {
	apiID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid api ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("api id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting api by ID", zap.Int("id", apiID))
	api, err := c.apiService.GetByID(apiID)
	if err != nil {
		c.Logger.Error("Error getting api by ID", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved api by ID", zap.Int("id", apiID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(api))
}

// UpdateApi
// @Summary update api
// @Description update api
// @Tags api
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseApi]
// @Router /v1/api [put]
func (c *ApiController) UpdateApi(ctx *gin.Context) {
	apiID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid api ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating api", zap.Int("id", apiID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for api update", zap.Error(err), zap.Int("id", apiID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for api update", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	apiUpdated, err := c.apiService.Update(apiID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating api", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseApi]().
		Data(domainToResponseMapper(apiUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Api updated successfully", zap.Int("id", apiID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteApi
// @Summary delete api
// @Description delete api by id
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/api/{id} [delete]
func (c *ApiController) DeleteApi(ctx *gin.Context) {
	apiID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid api ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting api", zap.Int("id", apiID))
	err = c.apiService.Delete(apiID)
	if err != nil {
		c.Logger.Error("Error deleting api", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Api deleted successfully", zap.Int("id", apiID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    apiID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// SearchApiPageList
// @Summary search apis
// @Description search apis by query
// @Tags search apis
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseApi]
// @Router /v1/api/search [get]
func (c *ApiController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching apis with pagination")

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
	for field := range apiRepo.ColumnsApiMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range apiRepo.ColumnsApiMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	fmt.Println(matches)
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range apiRepo.ColumnsApiMapping {
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

	result, err := c.apiService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching apis", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseApi]
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

	c.Logger.Info("Successfully searched apis",
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
// @Router /v1/api/search-property [get]
func (c *ApiController) SearchByProperty(ctx *gin.Context) {
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
		"apiName":   true,
		"email":     true,
		"firstName": true,
		"lastName":  true,
		"status":    true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.apiService.SearchByProperty(property, searchText)
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

// Mappers
func domainToResponseMapper(domainApi *domainApi.Api) *ResponseApi {

	return &ResponseApi{
		ID:          domainApi.ID,
		Path:        domainApi.Path,
		ApiGroup:    domainApi.ApiGroup,
		Method:      domainApi.Method,
		Description: domainApi.Description,
		CreatedAt:   domain.CustomTime{Time: domainApi.CreatedAt},
		UpdatedAt:   domain.CustomTime{Time: domainApi.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(apis *[]domainApi.Api) *[]*ResponseApi {
	res := make([]*ResponseApi, len(*apis))
	for i, u := range *apis {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewApiRequest) *domainApi.Api {
	return &domainApi.Api{
		Path:        req.Path,
		ApiGroup:    req.ApiGroup,
		Method:      req.Method,
		Description: req.Description,
	}
}
