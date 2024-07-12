package dto

import "time"

// OrderResponseDTO represents the response body for an order
type OrderResponseDTO struct {
	ID          uint                     `json:"id"`
	Bill        float64                  `json:"bill"`
	CurrentDate time.Time                `json:"currentDate"`
	Inventory   []InventoryResponseDTO   `json:"inventory"`
}

// InventoryResponseDTO represents the response body for inventory items
type InventoryResponseDTO struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}

// UserResponseDTO represents the structure of a user response
type UserResponseDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
