package routes

import (
	menuBtn "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menuBtn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuBtnRouters(router *gin.RouterGroup, controller menuBtn.IMenuBtnController) {
	u := router.Group("/menu_btn")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("", controller.NewMenuBtn)
		u.GET("", controller.GetAllMenuBtns)
		u.GET("/:id", controller.GetMenuBtnsByID)
		u.PUT("/:id", controller.UpdateMenuBtn)
		u.DELETE("/:id", controller.DeleteMenuBtn)
	}
}
