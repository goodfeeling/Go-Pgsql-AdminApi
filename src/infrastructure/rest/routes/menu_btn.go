package routes

import (
	"github.com/casbin/casbin/v2"
	menuBtn "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuBtnRouters(router *gin.RouterGroup, controller menuBtn.IMenuBtnController, enforcer *casbin.Enforcer) {
	u := router.Group("/menu_btn")
	u.Use(middlewares.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(enforcer))
	{
		u.POST("", controller.NewMenuBtn)
		u.GET("", controller.GetAllMenuBtns)
		u.GET("/:id", controller.GetMenuBtnsByID)
		u.PUT("/:id", controller.UpdateMenuBtn)
		u.DELETE("/:id", controller.DeleteMenuBtn)
	}
}
