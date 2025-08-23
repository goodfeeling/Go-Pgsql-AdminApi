package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/config"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func ConfigRouters(router *gin.RouterGroup, controller config.IConfigController, enforcer *casbin.Enforcer) {
	u := router.Group("/config")
	u.GET("/system", controller.GetConfigBySystem)

	u.Use(middlewares.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.GET("", controller.GetAllConfigs)
		u.PUT("/:module", controller.UpdateConfig)
		u.GET("/:module", controller.GetConfigByModule)
	}
}
