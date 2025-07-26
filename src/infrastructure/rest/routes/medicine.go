package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/medicine"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MedicineRoutes(router *gin.RouterGroup, controller medicine.IMedicineController, enforcer *casbin.Enforcer) {
	med := router.Group("/medicine")
	med.Use(middlewares.AuthJWTMiddleware())
	med.Use(middlewares.CasbinMiddleware(enforcer))
	{
		med.GET("/", controller.GetAllMedicines)
		med.POST("/", controller.NewMedicine)
		med.GET("/:id", controller.GetMedicinesByID)
		med.PUT("/:id", controller.UpdateMedicine)
		med.DELETE("/:id", controller.DeleteMedicine)
		med.GET("/search", controller.SearchPaginated)
		med.GET("/search-property", controller.SearchByProperty)
	}
}
