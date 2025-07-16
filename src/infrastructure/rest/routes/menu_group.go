package routes

import (
	menu_group "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menuGroup"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuGroupRouters(router *gin.RouterGroup, controller menu_group.IMenuGroupController) {
	u := router.Group("/menu_group")
	u.Use(middlewares.AuthJWTMiddleware())
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
