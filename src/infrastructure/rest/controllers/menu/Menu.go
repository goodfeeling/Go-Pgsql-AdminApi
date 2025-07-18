package menu

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenu "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
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
	MenuLevel   int               `json:"menu_level"`
	ParentID    int               `json:"parent_id"`
	Path        string            `json:"path"`
	Name        string            `json:"name"`
	Hidden      int16             `json:"hidden"`
	Component   string            `json:"component"`
	Sort        int8              `json:"sort"`
	ActiveName  string            `json:"active_name"`
	KeepAlive   int16             `json:"keep_alive"`
	DefaultMenu int16             `json:"default_menu"`
	Title       string            `json:"title"`
	Icon        string            `json:"icon"`
	CloseTab    int16             `json:"close_tab"`
	CreatedAt   domain.CustomTime `json:"created_at"`
	UpdatedAt   domain.CustomTime `json:"updated_at"`
}
type IMenuController interface {
	NewMenu(ctx *gin.Context)
	GetAllMenus(ctx *gin.Context)
	GetMenusByID(ctx *gin.Context)
	UpdateMenu(ctx *gin.Context)
	DeleteMenu(ctx *gin.Context)
	GetUserMenus(ctx *gin.Context)
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
	c.Logger.Info("Getting all menus by group id")
	groupId, err := strconv.Atoi(ctx.Query("group_id"))
	if err != nil {
		groupId = 0
	}
	menus, err := c.menuService.GetAll(groupId)
	if err != nil {
		c.Logger.Error("Error getting all menus", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	if menus == nil {
		menus = []*domainMenu.MenuTree{}
	}
	response := controllers.NewCommonResponseBuilder[[]*domainMenu.MenuTree]().
		Data(menus).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	ctx.JSON(http.StatusOK, response)
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

// GetUserMenus
// @Summary get user menus
// @Description user menus
// @Tags menus
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /v1/menu/user [get]
func (c *MenuController) GetUserMenus(ctx *gin.Context) {
	c.Logger.Info("Getting user menus")
	appUtils := controllers.NewAppUtils(ctx)
	userId, ok := appUtils.GetUserID()
	if !ok {
		c.Logger.Error("Error getting user id")
		appError := domainErrors.NewAppError(errors.New("user id not found"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menus, err := c.menuService.GetUserMenus(userId)
	if err != nil {
		c.Logger.Error("Error getting all menu user", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	if menus == nil {
		menus = []*domainMenu.MenuGroup{}
	}
	menuResponse := controllers.NewCommonResponseBuilder[[]*domainMenu.MenuGroup]().
		Data(menus).
		Message("success").
		Status(0).
		Build()
	ctx.JSON(http.StatusOK, menuResponse)

}

// Mappers
func domainToResponseMapper(domainMenu *domainMenu.Menu) *ResponseMenu {

	return &ResponseMenu{
		ID:          domainMenu.ID,
		Path:        domainMenu.Path,
		Name:        domainMenu.Name,
		ParentID:    domainMenu.ParentID,
		DefaultMenu: domainMenu.DefaultMenu,
		Hidden:      domainMenu.Hidden,
		MenuLevel:   domainMenu.MenuLevel,
		CloseTab:    domainMenu.CloseTab,
		KeepAlive:   domainMenu.KeepAlive,
		Icon:        domainMenu.Icon,
		Title:       domainMenu.Title,
		Sort:        domainMenu.Sort,
		ActiveName:  domainMenu.ActiveName,
		Component:   domainMenu.Component,
		CreatedAt:   domain.CustomTime{Time: domainMenu.CreatedAt},
		UpdatedAt:   domain.CustomTime{Time: domainMenu.UpdatedAt},
	}
}

func toUsecaseMapper(req *NewMenuRequest) *domainMenu.Menu {
	return &domainMenu.Menu{
		Path: req.Path,
	}
}
