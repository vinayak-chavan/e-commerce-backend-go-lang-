package controllers

import (
	"net/http"
	"time"

	"e-commerce/db"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// AddOrderFromCart creates an order from user's cart and stores it in Order and Inventory tables
func AddOrderFromCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	// Check if the user's cart has items
	var cartItems []models.Cart
	if err := db.DB.Where("user_id = ?", userIDUint).Preload("Product").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate total bill from cart items
	bill := calculateTotalBill(cartItems)

	// Begin transaction to ensure atomicity
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the order instance
	order := models.Order{
		UserId:      userIDUint,
		Bill:        bill,
		CurrentDate: time.Now(),
	}

	// Create the order in the database
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Add cart items to the inventory table with the order ID
	for _, cartItem := range cartItems {
		inventory := models.Inventory{
			OrderId:  order.ID,
			Name:     cartItem.Product.Name,
			Price:    cartItem.Product.Price,
			Quantity: cartItem.Quantity,
		}

		if err := tx.Create(&inventory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add inventory"})
			return
		}
	}

	// Clear user's cart after successful order creation
	if err := tx.Where("user_id = ?", userIDUint).Delete(&models.Cart{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Remove user details from the response for security reasons
	order.User = models.User{}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order": order})
}

// calculateTotalBill calculates the total bill from cart items
func calculateTotalBill(cartItems []models.Cart) float64 {
	var total float64
	for _, item := range cartItems {
		total += float64(item.Quantity) * item.Product.Price
	}
	return total
}

func GetMyOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	var orders []models.Order
	if err := db.DB.Preload("Inventory").Where("user_id = ?", userIDUint).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Prepare a clean response without nested Order and User details
	var cleanOrders []gin.H
	for _, order := range orders {
		var cleanInventory []gin.H
		for _, inv := range order.Inventory {
			cleanInventory = append(cleanInventory, gin.H{
				"id":       inv.ID,
				"name":     inv.Name,
				"price":    inv.Price,
				"quantity": inv.Quantity,
			})
		}

		cleanOrders = append(cleanOrders, gin.H{
			"id":          order.ID,
			"bill":        order.Bill,
			"currentDate": order.CurrentDate,
			"inventory":   cleanInventory,
		})
	}

	c.JSON(http.StatusOK, cleanOrders)
}

// GetAllOrders fetches all orders for admin role only
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	if err := db.DB.Preload("User").Preload("Inventory").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Prepare a clean response
	var response []gin.H
	for _, order := range orders {
		var cleanInventory []gin.H
		for _, inv := range order.Inventory {
			cleanInventory = append(cleanInventory, gin.H{
				"id":       inv.ID,
				"name":     inv.Name,
				"price":    inv.Price,
				"quantity": inv.Quantity,
			})
		}

		response = append(response, gin.H{
			"id":          order.ID,
			"bill":        order.Bill,
			"currentDate": order.CurrentDate,
			"user": gin.H{
				"id":    order.UserId,
				"name":  order.User.Name,
				"email": order.User.Email,
				// Add other user details as needed
			},
			"inventory": cleanInventory,
		})
	}

	c.JSON(http.StatusOK, response)
}
