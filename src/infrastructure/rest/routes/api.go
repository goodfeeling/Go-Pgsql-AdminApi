package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func ApiRouters(
	router *gin.RouterGroup,
	routerEngine *gin.Engine,
	controller api.IApiController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {

	// 用户获取接口列表
	if routerSetter, ok := controller.(api.RouterSetter); ok {
		routerSetter.SetRouter(routerEngine)
	}
	u := router.Group("/api")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewApi)
		u.GET("", controller.GetAllApis)
		u.GET("/:id", controller.GetApisByID)
		u.PUT("/:id", controller.UpdateApi)
		u.DELETE("/:id", controller.DeleteApi)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
		u.POST("/delete-batch", controller.DeleteApis)
		u.GET("/group-list", controller.GetApisGroup)
		u.POST("/synchronize", controller.SynchronizeRouterToApi)
	}
}
