package role

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainRole "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	roleRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewRoleRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	ParentID    int64  `json:"parent_id" binding:"required"`
	Order       int64  `json:"order"`
	Label       string `json:"label"`
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

type ResponseRole struct {
	ID            int64             `json:"id"`
	Name          string            `json:"name"`
	ParentID      int64             `json:"parent_id"`
	Order         int64             `json:"order"`
	Label         string            `json:"label"`
	Status        bool              `json:"status"`
	Description   string            `json:"description"`
	DefaultRouter string            `json:"default_router"`
	CreatedAt     domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt     domain.CustomTime `json:"updated_at,omitempty"`
}
type IRoleController interface {
	NewRole(ctx *gin.Context)
	GetAllRoles(ctx *gin.Context)
	GetRolesByID(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
}
type RoleController struct {
	roleService domainRole.IRoleService
	Logger      *logger.Logger
}

func NewRoleController(roleService domainRole.IRoleService, loggerInstance *logger.Logger) IRoleController {
	return &RoleController{roleService: roleService, Logger: loggerInstance}
}

func (c *RoleController) NewRole(ctx *gin.Context) {
	c.Logger.Info("Creating new role")
	var request NewRoleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new role", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	roleModel, err := c.roleService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating role", zap.Error(err), zap.String("Name", request.Name))
		_ = ctx.Error(err)
		return
	}
	roleResponse := controllers.NewCommonResponseBuilder[*ResponseRole]().
		Data(domainToResponseMapper(roleModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role created successfully", zap.String("Name", request.Name), zap.Int("id", int(roleModel.ID)))
	ctx.JSON(http.StatusOK, roleResponse)
}

func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	c.Logger.Info("Getting all roles")
	roles, err := c.roleService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all roles", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all roles", zap.Int("count", len(*roles)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]*ResponseRole]{
		Data: arrayDomainToResponseMapper(roles),
	})
}

func (c *RoleController) GetRolesByID(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("role id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting role by ID", zap.Int("id", roleID))
	role, err := c.roleService.GetByID(roleID)
	if err != nil {
		c.Logger.Error("Error getting role by ID", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved role by ID", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(role))
}

func (c *RoleController) UpdateRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating role", zap.Int("id", roleID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for role update", zap.Error(err), zap.Int("id", roleID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for role update", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	roleUpdated, err := c.roleService.Update(roleID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating role", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseRole]().
		Data(domainToResponseMapper(roleUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, response)
}

func (c *RoleController) DeleteRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting role", zap.Int("id", roleID))
	err = c.roleService.Delete(roleID)
	if err != nil {
		c.Logger.Error("Error deleting role", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Role deleted successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully", "status": 0, "data": roleID})
}

func (c *RoleController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching roles with pagination")

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
	for field := range roleRepo.ColumnsRoleMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range roleRepo.ColumnsRoleMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range roleRepo.ColumnsRoleMapping {
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

	result, err := c.roleService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching roles", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseRole]
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

	c.Logger.Info("Successfully searched roles",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

func (c *RoleController) SearchByProperty(ctx *gin.Context) {
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
		"roleName":  true,
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

	coincidences, err := c.roleService.SearchByProperty(property, searchText)
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
func domainToResponseMapper(domainRole *domainRole.Role) *ResponseRole {
	return &ResponseRole{
		ID:          domainRole.ID,
		Name:        domainRole.Name,
		ParentID:    domainRole.ParentID,
		Order:       domainRole.Order,
		Label:       domainRole.Label,
		Description: domainRole.Description,
		Status:      domainRole.Status,
		CreatedAt:   domain.CustomTime{Time: domainRole.CreatedAt},
		UpdatedAt:   domain.CustomTime{Time: domainRole.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(roles *[]domainRole.Role) *[]*ResponseRole {
	res := make([]*ResponseRole, len(*roles))
	for i, u := range *roles {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewRoleRequest) *domainRole.Role {
	return &domainRole.Role{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		Order:       req.Order,
		Label:       req.Label,
		Status:      req.Status,
	}
}
