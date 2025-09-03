package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/scheduled_task"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func ScheduledTaskRouters(
	router *gin.RouterGroup,
	controller scheduled_task.IScheduledTaskController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/scheduled_task")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewScheduledTask)
		u.GET("", controller.GetAllScheduledTasks)
		u.GET("/:id", controller.GetScheduledTaskByID)
		u.PUT("/:id", controller.UpdateScheduledTask)
		u.DELETE("/:id", controller.DeleteScheduledTask)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
		u.POST("/delete-batch", controller.DeleteScheduledTasks)
		u.POST("/enable/:id", controller.EnableTaskById)
		u.POST("/disable/:id", controller.DisableTaskById)
		u.POST("/reload", controller.ReloadAllTasks)
	}
}
