package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func WebSocketRoute(router *gin.Engine, appContext *di.ApplicationContext) {
	appContext.WsRouter.AddRoute("/scheduleLog", appContext.TaskExecutionLogModule.WsHandler)

	r := router.Group("/ws")
	r.Use(middlewares.UrlAuthMiddleware())
	r.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	r.GET("scheduleLog", func(ctx *gin.Context) {
		appContext.WsRouter.HandleConnectionWithRoute(ctx, "/scheduleLog")
	})
}
