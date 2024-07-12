package controllers

import (
	"net/http"
	"time"

	"e-commerce/db"
	"e-commerce/dto"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// AddOrderFromCart creates an order from user's cart and stores it in Order and Inventory tables
// @Summary Add an order from the cart
// @Description Create an order from the user's cart
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} dto.OrderResponseDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /orders [post]
func AddOrderFromCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	// Check if the user's cart has items
	var cartItems []models.Cart
	if err := db.DB.Where("user_id = ?", userIDUint).Preload("Product").Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch cart items"})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Cart is empty"})
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create order"})
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
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to add inventory"})
			return
		}
	}

	// Clear user's cart after successful order creation
	if err := tx.Where("user_id = ?", userIDUint).Delete(&models.Cart{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to clear cart"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Prepare the order response DTO
	orderResponse := dto.OrderResponseDTO{
		ID:          order.ID,
		Bill:        order.Bill,
		CurrentDate: order.CurrentDate,
		Inventory:   mapToInventoryDTOs(order.Inventory),
	}

	c.JSON(http.StatusCreated, orderResponse)
}

// calculateTotalBill calculates the total bill from cart items
func calculateTotalBill(cartItems []models.Cart) float64 {
	var total float64
	for _, item := range cartItems {
		total += float64(item.Quantity) * item.Product.Price
	}
	return total
}

// mapToInventoryDTOs maps inventory items to inventory response DTOs
func mapToInventoryDTOs(inventoryItems []models.Inventory) []dto.InventoryResponseDTO {
	var inventoryDTOs []dto.InventoryResponseDTO
	for _, item := range inventoryItems {
		inventoryDTOs = append(inventoryDTOs, dto.InventoryResponseDTO{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}
	return inventoryDTOs
}

// GetMyOrders retrieves the user's orders
// @Summary Get my orders
// @Description Retrieve the current user's orders
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dto.OrderResponseDTO
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /orders [get]
func GetMyOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDUint, _ := userID.(uint)

	var orders []models.Order
	if err := db.DB.Preload("Inventory").Where("user_id = ?", userIDUint).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch orders"})
		return
	}

	// Prepare order response DTOs
	var orderResponses []dto.OrderResponseDTO
	for _, order := range orders {
		orderResponses = append(orderResponses, dto.OrderResponseDTO{
			ID:          order.ID,
			Bill:        order.Bill,
			CurrentDate: order.CurrentDate,
			Inventory:   mapToInventoryDTOs(order.Inventory),
		})
	}

	c.JSON(http.StatusOK, orderResponses)
}

// GetAllOrders fetches all orders for admin role only
// @Summary Get all orders
// @Description Retrieve all orders (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dto.OrderResponseDTO
// @Failure 500 {object} dto.ErrorResponse
// @securityDefinitions.apiKey Authorization
// @in header
// @name Authorization
// @Security JWT
// @Router /orders/all [get]
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	if err := db.DB.Preload("User").Preload("Inventory").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Prepare order response DTOs
	var orderResponses []dto.OrderResponseDTO
	for _, order := range orders {
		orderResponses = append(orderResponses, dto.OrderResponseDTO{
			ID:          order.ID,
			Bill:        order.Bill,
			CurrentDate: order.CurrentDate,
			Inventory:   mapToInventoryDTOs(order.Inventory),
		})
	}

	c.JSON(http.StatusOK, orderResponses)
}
