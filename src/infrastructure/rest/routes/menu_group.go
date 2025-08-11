package routes

import (
	"github.com/casbin/casbin/v2"
	menu_group "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu_group"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuGroupRouters(router *gin.RouterGroup, controller menu_group.IMenuGroupController, enforcer *casbin.Enforcer) {
	u := router.Group("/menu_group")
	u.Use(middlewares.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewMenuGroup)
		u.GET("", controller.GetAllMenuGroups)
		u.GET("/:id", controller.GetMenuGroupsByID)
		u.PUT("/:id", controller.UpdateMenuGroup)
		u.DELETE("/:id", controller.DeleteMenuGroup)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
