package routes

import (
	"github.com/casbin/casbin/v2"
	menuParameter "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu_parameter"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuParameterRouters(
	router *gin.RouterGroup,
	controller menuParameter.IMenuParameterController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/menu_parameter")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewMenuParameter)
		u.GET("", controller.GetAllMenuParameters)
		u.GET("/:id", controller.GetMenuParametersByID)
		u.PUT("/:id", controller.UpdateMenuParameter)
		u.DELETE("/:id", controller.DeleteMenuParameter)
	}
}
