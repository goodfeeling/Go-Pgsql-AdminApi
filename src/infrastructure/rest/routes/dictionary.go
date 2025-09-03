package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionary"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func DictionaryRouters(
	router *gin.RouterGroup,
	controller dictionary.IDictionaryController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/dictionary")

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
		u.GET("/type/:type", controller.GetByType)

	}
}
