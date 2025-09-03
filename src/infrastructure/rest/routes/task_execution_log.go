package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func TaskExecutionLogRouters(
	router *gin.RouterGroup,
	controller task_execution_log.ITaskExecutionLogController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/task_execution_log")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.GET("/search", controller.SearchPaginated)
	}
}
