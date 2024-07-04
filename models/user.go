package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Carts    []Cart  `gorm:"foreignKey:UserId"`
	Orders   []Order `gorm:"foreignKey:UserId"`
}
