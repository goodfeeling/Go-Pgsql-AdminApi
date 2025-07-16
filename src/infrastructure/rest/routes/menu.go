package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuRouters(router *gin.RouterGroup, controller menu.IMenuController) {
	u := router.Group("/menu")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("", controller.NewMenu)
		u.GET("", controller.GetAllMenus)
		u.GET("/:id", controller.GetMenusByID)
		u.PUT("/:id", controller.UpdateMenu)
		u.DELETE("/:id", controller.DeleteMenu)
		u.GET("/tree", controller.GetTreeMenus)
		u.GET("/user", controller.GetUserMenus)
	}
}
