package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.Engine) {
	productRoutes := router.Group("/products")
	{
		productRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetProducts)
		productRoutes.POST("/", middlewares.AuthMiddleware(), controllers.CreateProduct)
		productRoutes.GET("/:id", middlewares.AuthMiddleware(), controllers.GetProductByID)
		productRoutes.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateProduct)
		productRoutes.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteProduct)
	}
}
