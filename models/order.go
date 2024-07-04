package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId      uint        `json:"userId"`
	User        User        `gorm:"foreignKey:UserId"`
	Bill        float64     `json:"bill"`
	CurrentDate time.Time   `json:"currentDate"`
	Inventory   []Inventory `gorm:"foreignKey:OrderId"`
}
