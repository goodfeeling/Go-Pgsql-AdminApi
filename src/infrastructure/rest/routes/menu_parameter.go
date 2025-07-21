package routes

import (
	menuParameter "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menuParameter"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuParameterRouters(router *gin.RouterGroup, controller menuParameter.IMenuParameterController) {
	u := router.Group("/menu_parameter")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("", controller.NewMenuParameter)
		u.GET("", controller.GetAllMenuParameters)
		u.GET("/:id", controller.GetMenuParametersByID)
		u.PUT("/:id", controller.UpdateMenuParameter)
		u.DELETE("/:id", controller.DeleteMenuParameter)
	}
}
