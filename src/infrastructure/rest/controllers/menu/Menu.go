package menu

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenu "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewMenuRequest struct {
	ID          int    `json:"id"`
	Path        string `json:"path" binding:"required"`
	MenuGroup   string `json:"menu_group" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ResponseMenu struct {
	ID          int               `json:"id"`
	Path        string            `json:"path"`
	MenuGroup   string            `json:"menu_group"`
	Method      string            `json:"method"`
	Description string            `json:"description"`
	CreatedAt   domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt   domain.CustomTime `json:"updated_at,omitempty"`
}
type IMenuController interface {
	NewMenu(ctx *gin.Context)
	GetAllMenus(ctx *gin.Context)
	GetMenusByID(ctx *gin.Context)
	UpdateMenu(ctx *gin.Context)
	DeleteMenu(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
}
type MenuController struct {
	menuService domainMenu.IMenuService
	Logger      *logger.Logger
}

func NewMenuController(menuService domainMenu.IMenuService, loggerInstance *logger.Logger) IMenuController {
	return &MenuController{menuService: menuService, Logger: loggerInstance}
}

// CreateMenu
// @Summary create menu
// @Description create menu
// @Tags menu create
// @Accept json
// @Produce json
// @Param book body NewMenuRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/menu [post]
func (c *MenuController) NewMenu(ctx *gin.Context) {
	c.Logger.Info("Creating new menu")
	var request NewMenuRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new menu", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menuModel, err := c.menuService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating menu", zap.Error(err), zap.String("path", request.Path))
		_ = ctx.Error(err)
		return
	}
	menuResponse := controllers.NewCommonResponseBuilder[*ResponseMenu]().
		Data(domainToResponseMapper(menuModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Menu created successfully", zap.String("path", request.Path), zap.Int("id", int(menuModel.ID)))
	ctx.JSON(http.StatusOK, menuResponse)
}

// GetAllMenus
// @Summary get all menus by
// @Description get  all menus by where
// @Tags menus
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainMenu.Menu]
// @Router /v1/menu [get]
func (c *MenuController) GetAllMenus(ctx *gin.Context) {
	c.Logger.Info("Getting all menus")
	menus, err := c.menuService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all menus", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(*menus)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainMenu.Menu]{
		Data: menus,
	})
}

// GetMenusByID
// @Summary get menus
// @Description get menus by id
// @Tags menus
// @Accept json
// @Produce json
// @Success 200 {object} ResponseMenu
// @Router /v1/menu/{id} [get]
func (c *MenuController) GetMenusByID(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menu ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("menu id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting menu by ID", zap.Int("id", menuID))
	menu, err := c.menuService.GetByID(menuID)
	if err != nil {
		c.Logger.Error("Error getting menu by ID", zap.Error(err), zap.Int("id", menuID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved menu by ID", zap.Int("id", menuID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(menu))
}

// UpdateMenu
// @Summary update menu
// @Description update menu
// @Tags menu
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseMenu]
// @Router /v1/menu [put]
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menu ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating menu", zap.Int("id", menuID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for menu update", zap.Error(err), zap.Int("id", menuID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for menu update", zap.Error(err), zap.Int("id", menuID))
		_ = ctx.Error(err)
		return
	}
	menuUpdated, err := c.menuService.Update(menuID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating menu", zap.Error(err), zap.Int("id", menuID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseMenu]().
		Data(domainToResponseMapper(menuUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Menu updated successfully", zap.Int("id", menuID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteMenu
// @Summary delete menu
// @Description delete menu by id
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/menu/{id} [delete]
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menu ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting menu", zap.Int("id", menuID))
	err = c.menuService.Delete(menuID)
	if err != nil {
		c.Logger.Error("Error deleting menu", zap.Error(err), zap.Int("id", menuID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Menu deleted successfully", zap.Int("id", menuID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    menuID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// SearchMenuPageList
// @Summary search menus
// @Description search menus by query
// @Tags search menus
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseMenu]
// @Router /v1/menu/search [get]
func (c *MenuController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching menus with pagination")

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
	for field := range menuRepo.ColumnsMenuMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range menuRepo.ColumnsMenuMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	fmt.Println(matches)
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range menuRepo.ColumnsMenuMapping {
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

	result, err := c.menuService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching menus", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseMenu]
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

	c.Logger.Info("Successfully searched menus",
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
// @Router /v1/menu/search-property [get]
func (c *MenuController) SearchByProperty(ctx *gin.Context) {
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
		"menuName":  true,
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

	coincidences, err := c.menuService.SearchByProperty(property, searchText)
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
func domainToResponseMapper(domainMenu *domainMenu.Menu) *ResponseMenu {

	return &ResponseMenu{
		ID:          domainMenu.ID,
		Path:        domainMenu.Path,
		MenuGroup:   domainMenu.MenuGroup,
		Method:      domainMenu.Method,
		Description: domainMenu.Description,
		CreatedAt:   domain.CustomTime{Time: domainMenu.CreatedAt},
		UpdatedAt:   domain.CustomTime{Time: domainMenu.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(menus *[]domainMenu.Menu) *[]*ResponseMenu {
	res := make([]*ResponseMenu, len(*menus))
	for i, u := range *menus {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewMenuRequest) *domainMenu.Menu {
	return &domainMenu.Menu{
		Path:        req.Path,
		MenuGroup:   req.MenuGroup,
		Method:      req.Method,
		Description: req.Description,
	}
}
