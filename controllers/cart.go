package controllers

import (
	"net/http"
	"strconv"

	"e-commerce/db"
	"e-commerce/dto"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// AddToCart adds a product to the user's cart
// @Summary Add a product to the cart
// @Description Add a specific product to the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param AddToCartRequest body dto.AddToCartRequest true "Add to Cart Request"
// @Success 200 {object} dto.CartItemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /cart [post]
func AddToCart(c *gin.Context) {
	var input struct {
		ProductID uint `json:"productId" binding:"required"`
		Quantity  uint `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	// Check if the product exists
	var product models.Product
	if err := db.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
		return
	}

	// Create or update the cart item
	var cartItem models.Cart
	if err := db.DB.Where("user_id = ? AND product_id = ?", userIDUint, input.ProductID).First(&cartItem).Error; err != nil {
		// Create new cart item
		cartItem = models.Cart{
			UserId:    userIDUint,
			ProductId: input.ProductID,
			Quantity:  input.Quantity,
		}
		if err := db.DB.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to add item to cart"})
			return
		}
	} else {
		// Update existing cart item
		cartItem.Quantity += input.Quantity
		if err := db.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update cart item"})
			return
		}
	}

	c.JSON(http.StatusOK, cartItem)
}

// ViewCart retrieves the user's cart items
// @Summary View cart items
// @Description Retrieve all items in the user's cart
// @Tags cart
// @Produce json
// @Success 200 {array} dto.CartItemResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /cart [get]
func ViewCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	var cartItems []models.Cart

	// Fetch cart items with preloaded product details
	if err := db.DB.Preload("Product").Where("user_id = ?", userIDUint).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch cart items"})
		return
	}

	// Define a slice to hold the formatted response
	var cartResponse []gin.H

	// Iterate through each cart item and format the response
	for _, cart := range cartItems {
		cartResponse = append(cartResponse, gin.H{
			"id":        cart.ID,
			"productId": cart.ProductId,
			"product": gin.H{
				"id":          cart.Product.ID,
				"name":        cart.Product.Name,
				"description": cart.Product.Description,
				"price":       cart.Product.Price,
				"photo":       cart.Product.Photo,
			},
			"quantity": cart.Quantity,
		})
	}

	c.JSON(http.StatusOK, cartResponse)
}

// RemoveFromCart removes an item from the user's cart
// @Summary Remove an item from the cart
// @Description Remove a specific item from the user's cart by cart item ID
// @Tags cart
// @Produce json
// @Param id path uint true "Cart Item ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /cart/{id} [delete]
func RemoveFromCart(c *gin.Context) {
	id := c.Param("id")
	cartID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid cart item ID"})
		return
	}

	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	// Check if the cart item exists
	var cartItem models.Cart
	if err := db.DB.Where("id = ? AND user_id = ?", cartID, userIDUint).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Cart item not found"})
		return
	}

	// Delete the cart item
	if err := db.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Item removed from cart"})

}

// UpdateCartItem updates the quantity of a specific item in the user's cart
// @Summary Update cart item quantity
// @Description Update the quantity of a specific item in the user's cart by cart item ID
// @Tags cart
// @Accept json
// @Produce json
// @Param id path uint true "Cart Item ID"
// @Param UpdateQuantity body dto.UpdateQuantity true "Change quantity"
// @Success 200 {object} dto.CartItemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /cart/{id} [put]
func UpdateCartItem(c *gin.Context) {
	var input struct {
		Quantity uint `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	id := c.Param("id")
	cartID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid cart item ID"})
		return
	}

	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	// Check if the cart item exists
	var cartItem models.Cart
	if err := db.DB.Where("id = ? AND user_id = ?", cartID, userIDUint).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Cart item not found"})
		return
	}

	// Update the quantity
	cartItem.Quantity = input.Quantity
	if err := db.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}
