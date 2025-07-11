package routes

import (
	dictionary_detail "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionaryDetail"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func DictionaryDetailRouters(router *gin.RouterGroup, controller dictionary_detail.IIDictionaryDetailController) {
	u := router.Group("/dictionary_detail")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("", controller.NewDictionary)
		u.GET("", controller.GetAllDictionaries)
		u.GET("/:id", controller.GetDictionariesByID)
		u.PUT("/:id", controller.UpdateDictionary)
		u.DELETE("/:id", controller.DeleteDictionary)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
		u.POST("/delete-batch", controller.DeleteDictionaryDetails)
	}
}
