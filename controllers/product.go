package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"e-commerce/db"
	"e-commerce/dto"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// GetProducts fetches all products
// @Summary Get all products
// @Description Retrieve a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// CreateProduct creates a new product with photo upload
// @Summary Create a new product
// @Description Create a new product with the given details
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param description formData string true "Product Description"
// @Param price formData number true "Product Price"
// @Param photo formData file true "Product Photo"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Retrieve file from the request
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Photo is required"})
		return
	}

	// Generate the file path
	timestamp := time.Now().Unix()
	filePath := filepath.Join("uploads", fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename)))

	// Save the file to the specified path
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Unable to save the photo"})
		return
	}

	// Bind the other product details
	product.Name = c.PostForm("name")
	product.Description = c.PostForm("description")
	price := c.PostForm("price")
	fmt.Sscanf(price, "%f", &product.Price)
	product.Photo = filePath

	if err := db.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
// @Summary Update an existing product
// @Description Update an existing product with the given details
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param name formData string false "Product Name"
// @Param description formData string false "Product Description"
// @Param price formData number false "Product Price"
// @Param photo formData file false "Product Photo"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if the product exists
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
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
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Unable to save the photo"})
			return
		}

		product.Photo = filePath // Update photo path in the product struct
	}

	// Bind the other product details from form data
	var updateData dto.ProductRequest

	if err := c.ShouldBind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Respond with updated product details
	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Check if the product exists
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
		return
	}

	// Delete the product
	if err := db.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Product deleted successfully"})
}

// GetProductByID fetches details of a product by ID
// @Summary Get a product by ID
// @Description Retrieve details of a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 404 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Retrieve the product by ID
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
