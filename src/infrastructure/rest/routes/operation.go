package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/operation"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func OperationRouters(router *gin.RouterGroup, controller operation.IOperationController, enforcer *casbin.Enforcer) {
	u := router.Group("/operation")
	u.Use(middlewares.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.GET("", controller.GetAllOperations)
		u.GET("/:id", controller.GetOperationsByID)
		u.DELETE("/:id", controller.DeleteOperation)
		u.POST("/delete-batch", controller.DeleteOperations)
		u.GET("/search", controller.SearchPaginated)
	}
}
