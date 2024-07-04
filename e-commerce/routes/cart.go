package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(router *gin.Engine) {
	cartRoutes := router.Group("/cart")
	{
		cartRoutes.Use(middlewares.AuthMiddleware())
		cartRoutes.POST("/", controllers.AddToCart)
		cartRoutes.GET("/", controllers.ViewCart)
		cartRoutes.PUT("/:id", controllers.UpdateCartItem)
		cartRoutes.DELETE("/:id", controllers.RemoveFromCart)
	}
}
