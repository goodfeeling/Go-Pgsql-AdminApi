package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/upload"
	"github.com/gin-gonic/gin"
)

func UploadRoutes(router *gin.RouterGroup, controller upload.IUploadController) {
	u := router.Group("/upload")
	{
		u.POST("/single", controller.Single)
		u.POST("/multiple", controller.Multiple)
	}
}
