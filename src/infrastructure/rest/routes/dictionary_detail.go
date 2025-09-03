package routes

import (
	"github.com/casbin/casbin/v2"
	dictionary_detail "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionary_detail"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func DictionaryDetailRouters(
	router *gin.RouterGroup,
	controller dictionary_detail.IIDictionaryDetailController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/dictionary_detail")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
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
