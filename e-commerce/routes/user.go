package routes

import (
    "e-commerce/controllers"
    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
    userRoutes := router.Group("/users")
    {
        userRoutes.POST("/register", controllers.RegisterUser)
        userRoutes.POST("/login", controllers.LoginUser)
    }
}