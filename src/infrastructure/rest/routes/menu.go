package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuRouters(router *gin.RouterGroup, controller menu.IMenuController, enforcer *casbin.Enforcer) {
	router.Use(middlewares.OptionalAuthMiddleware()).GET("/menu/user", controller.GetUserMenus)
	u := router.Group("/menu")
	u.Use(middlewares.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewMenu)
		u.GET("", controller.GetAllMenus)
		u.GET("/:id", controller.GetMenusByID)
		u.PUT("/:id", controller.UpdateMenu)
		u.DELETE("/:id", controller.DeleteMenu)
	}
}
