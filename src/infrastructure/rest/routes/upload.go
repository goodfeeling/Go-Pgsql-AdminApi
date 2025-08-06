package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/upload"
	"github.com/gin-gonic/gin"
)

func UploadRoutes(router *gin.RouterGroup, controller upload.IUploadController, enforcer *casbin.Enforcer) {
	u := router.Group("/upload")
	// u.Use(middlewares.AuthJWTMiddleware())
	// u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("/single", controller.Single)
		u.POST("/multiple", controller.Multiple)
		u.GET("/sts-token", controller.GetSTSToken)
	}
}
