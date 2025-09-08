package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.MenuModule.Controller
	middlewareProvider := appContext.MiddlewareProvider
	router.Use(middlewareProvider.OptionalAuthMiddleware()).GET("/menu/user", controller.GetUserMenus)

	u := router.Group("/menu")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewMenu)
		u.GET("", controller.GetAllMenus)
		u.GET("/:id", controller.GetMenusByID)
		u.PUT("/:id", controller.UpdateMenu)
		u.DELETE("/:id", controller.DeleteMenu)
	}
}
