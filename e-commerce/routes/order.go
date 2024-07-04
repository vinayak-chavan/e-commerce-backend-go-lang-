package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.Engine) {
	productRoutes := router.Group("/orders")
	{
		productRoutes.POST("/", middlewares.AuthMiddleware(), controllers.AddOrderFromCart)
		productRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetMyOrders)
		productRoutes.GET("/all", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.GetAllOrders)
	}
}
