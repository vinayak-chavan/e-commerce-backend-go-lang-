package models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	OrderId  uint    `json:"orderId"`
	Order    Order   `gorm:"foreignKey:OrderId"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}
