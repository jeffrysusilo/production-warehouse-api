package routes

import (
	"github.com/gin-gonic/gin"
	"production-warehouse-api/controllers"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/items", controllers.CreateItem)
	router.GET("/items", controllers.GetItems)
	router.GET("/items/:id", controllers.GetItemByID)
	router.PUT("/items/:id", controllers.UpdateItem)
	router.DELETE("/items/:id", controllers.DeleteItem)
}
