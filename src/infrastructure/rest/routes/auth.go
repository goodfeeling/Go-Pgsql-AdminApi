package routes

import (
	"github.com/casbin/casbin/v2"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, controller authController.IAuthController, enforcer *casbin.Enforcer,
	middlewareProvider *middlewares.MiddlewareProvider) {
	routerAuth := router.Group("/auth")
	{
		routerAuth.POST("/signin", controller.Login)
		routerAuth.POST("/signup", controller.Register)
		routerAuth.POST("/access-token", controller.GetAccessTokenByRefreshToken)
	}
	loginAuth := routerAuth.Use(middlewareProvider.AuthJWTMiddleware())
	{
		loginAuth.POST("/switch-role", controller.SwitchRole)
		loginAuth.GET("/logout", controller.Logout)
	}
}
