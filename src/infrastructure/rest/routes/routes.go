package routes

import (
	"net/http"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	AuthRoutes(v1, appContext.AuthModule.Controller, appContext.Enforcer)
	UserRoutes(v1, appContext.UserModule.Controller, appContext.Enforcer)
	UploadRoutes(v1, appContext.UploadModule.Controller, appContext.Enforcer)
	RoleRoutes(v1, appContext.RoleModule.Controller, appContext.Enforcer)
	ApiRouters(v1, router, appContext.ApiModule.Controller, appContext.Enforcer)
	OperationRouters(v1, appContext.OperationModule.Controller, appContext.Enforcer)
	DictionaryRouters(v1, appContext.DictionaryModule.Controller, appContext.Enforcer)
	DictionaryDetailRouters(v1, appContext.DictionaryDetailModule.Controller, appContext.Enforcer)
	MenuRouters(v1, appContext.MenuModule.Controller, appContext.Enforcer)
	MenuGroupRouters(v1, appContext.MenuGroupModule.Controller, appContext.Enforcer)
	MenuBtnRouters(v1, appContext.MenuBtnModule.Controller, appContext.Enforcer)
	MenuParameterRouters(v1, appContext.MenuParameterModule.Controller, appContext.Enforcer)
	FileRouters(v1, appContext.FileModule.Controller, appContext.Enforcer)

	ScheduledTaskRouters(v1, appContext.ScheduledTaskModule.Controller, appContext.Enforcer)
	TaskExecutionLogRouters(v1, appContext.TaskExecutionLogModule.Controller, appContext.Enforcer)
	ConfigRouters(v1, appContext.ConfigModule.Controller, appContext.Enforcer)
}
