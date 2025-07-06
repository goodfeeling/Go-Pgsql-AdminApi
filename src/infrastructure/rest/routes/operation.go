package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func OperationRouters(router *gin.RouterGroup, controller api.IApiController) {
	u := router.Group("/api")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("", controller.NewApi)
		u.GET("", controller.GetAllApis)
		u.GET("/:id", controller.GetApisByID)
		u.PUT("/:id", controller.UpdateApi)
		u.DELETE("/:id", controller.DeleteApi)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
