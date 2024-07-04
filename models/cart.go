package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId    uint    `json:"userId"`
	User      User    `gorm:"foreignKey:UserId"`
	ProductId uint    `json:"productId"`
	Product   Product `gorm:"foreignKey:ProductId"`
	Quantity  uint    `json:"quantity"`
}
