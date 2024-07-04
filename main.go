package main

import (
	"log"
	"os"

	"e-commerce/db"
	"e-commerce/models"
	"e-commerce/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	db.InitDatabase()
	models.MigrateDatabase()

	router := gin.Default()

	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Application working fine!!",
		})
	})

	os.MkdirAll("uploads", os.ModePerm)
	router.Static("/uploads", "./uploads")

	routes.RegisterProductRoutes(router)
	routes.RegisterUserRoutes(router)
	routes.RegisterCartRoutes(router)
	routes.RegisterOrderRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router.Run(":" + port)
}
