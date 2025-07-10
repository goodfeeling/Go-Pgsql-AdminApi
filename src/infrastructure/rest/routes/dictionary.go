package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionary"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func DictionaryRouters(router *gin.RouterGroup, controller dictionary.IDictionaryController) {
	u := router.Group("/dictionary")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("", controller.NewDictionary)
		u.GET("", controller.GetAllDictionaries)
		u.GET("/:id", controller.GetDictionariesByID)
		u.PUT("/:id", controller.UpdateDictionary)
		u.DELETE("/:id", controller.DeleteDictionary)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
