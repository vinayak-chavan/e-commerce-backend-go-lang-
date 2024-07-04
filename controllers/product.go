package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"e-commerce/db"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// GetProducts fetches all products
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// CreateProduct creates a new product with photo upload
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Retrieve file from the request
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo is required"})
		return
	}

	// Generate the file path
	timestamp := time.Now().Unix()
	filePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename)))

	// Save the file to the specified path
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the photo"})
		return
	}

	// Bind the other product details
	product.Name = c.PostForm("name")
	product.Description = c.PostForm("description")
	price := c.PostForm("price")
	fmt.Sscanf(price, "%f", &product.Price)
	product.Photo = filePath

	if err := db.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if the product exists
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Retrieve file from the request, if any
	file, err := c.FormFile("photo")
	if err == nil {
		// Generate the file path
		timestamp := time.Now().Unix()
		filePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename)))

		// Save the file to the specified path
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the photo"})
			return
		}

		product.Photo = filePath // Update photo path in the product struct
	}

	// Bind the other product details from form data
	var updateData struct {
		Name        string  `form:"name"`
		Description string  `form:"description"`
		Price       float64 `form:"price"`
	}

	if err := c.ShouldBind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields only if they are provided
	if updateData.Name != "" {
		product.Name = updateData.Name
	}
	if updateData.Description != "" {
		product.Description = updateData.Description
	}
	if updateData.Price != 0 {
		product.Price = updateData.Price
	}

	// Update the product in the database
	if err := db.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with updated product details
	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if the product exists
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Delete the product
	if err := db.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetProductByID fetches details of a product by ID
func GetProductByID(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Retrieve the product by ID
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
