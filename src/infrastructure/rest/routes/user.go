package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(
	router *gin.RouterGroup,
	controller user.IUserController,
	enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	u := router.Group("/user")
	u.Use(middlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewUser)
		u.GET("", controller.GetAllUsers)
		u.GET("/:id", controller.GetUsersByID)
		u.PUT("/:id", controller.UpdateUser)
		u.DELETE("/:id", controller.DeleteUser)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
		u.POST(":id/role", controller.UserBindRoles)
		u.POST("/:id/reset-password", controller.ResetPassword)
		u.POST("/:id/edit-password", controller.EditPassword)
	}
}
