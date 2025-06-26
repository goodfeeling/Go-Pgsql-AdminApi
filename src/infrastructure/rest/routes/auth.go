package routes

import (
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, controller authController.IAuthController) {
	routerAuth := router.Group("/auth")
	{
		routerAuth.POST("/signin", controller.Login)
		routerAuth.POST("/signup", controller.Register)
		routerAuth.POST("/access-token", controller.GetAccessTokenByRefreshToken)
		routerAuth.Use(middlewares.AuthJWTMiddleware()).GET("/logout", controller.Logout)
	}
}
