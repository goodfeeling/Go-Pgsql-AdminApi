package routes

import (
	"net/http"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {
	v1 := router.Group("/v1")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	AuthRoutes(v1, appContext.AuthModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	UserRoutes(v1, appContext.UserModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	UploadRoutes(v1, appContext.UploadModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	RoleRoutes(v1, appContext.RoleModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	ApiRouters(v1, router, appContext.ApiModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	OperationRouters(v1, appContext.OperationModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	DictionaryRouters(v1, appContext.DictionaryModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	DictionaryDetailRouters(v1, appContext.DictionaryDetailModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	MenuRouters(v1, appContext.MenuModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	MenuGroupRouters(v1, appContext.MenuGroupModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	MenuBtnRouters(v1, appContext.MenuBtnModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	MenuParameterRouters(v1, appContext.MenuParameterModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	FileRouters(v1, appContext.FileModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)

	ScheduledTaskRouters(v1, appContext.ScheduledTaskModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	ConfigRouters(v1, appContext.ConfigModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)
	TaskExecutionLogRouters(v1, appContext.TaskExecutionLogModule.Controller, appContext.Enforcer, appContext.MiddlewareProvider)

}
