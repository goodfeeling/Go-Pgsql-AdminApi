package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(
	router *gin.RouterGroup,
	controller role.IRoleController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/role")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewRole)
		u.GET("", controller.GetAllRoles)
		u.GET("/:id", controller.GetRolesByID)
		u.PUT("/:id", controller.UpdateRole)
		u.DELETE("/:id", controller.DeleteRole)
		u.GET("/tree", controller.GetTreeRoles)
		u.GET("/:id/setting", controller.GetRoleSetting)
		u.POST("/:id/menu", controller.UpdateRoleMenuIds)
		u.POST("/:id/api", controller.BindApiRule)
		u.POST("/:id/menu-btns", controller.BindRoleMenuBtns)
	}
}
