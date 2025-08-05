package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/file"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func FileRouters(router *gin.RouterGroup, controller file.IFileController, enforcer *casbin.Enforcer) {
	u := router.Group("/file")
	u.Use(middlewares.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewFile)
		u.GET("", controller.GetAllFiles)
		u.GET("/:id", controller.GetFilesByID)
		u.PUT("/:id", controller.UpdateFile)
		u.DELETE("/:id", controller.DeleteFile)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
