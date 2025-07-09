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

	AuthRoutes(v1, appContext.AuthController)
	UserRoutes(v1, appContext.UserController)
	MedicineRoutes(v1, appContext.MedicineController)
	UploadRoutes(v1, appContext.UploadController)
	RoleRoutes(v1, appContext.RoleController)
	ApiRouters(v1, appContext.ApiController)
	OperationRouters(v1, appContext.OperationController)
	DictionaryRouters(v1, appContext.DictionaryController)
	DictionaryDetailRouters(v1, appContext.DictionaryDetailController)
}
