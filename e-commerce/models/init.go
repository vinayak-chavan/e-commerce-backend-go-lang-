package models

import (
	"e-commerce/db"
	"log"
)

func MigrateDatabase() {
	err := db.DB.AutoMigrate(&User{}, &Product{}, &Cart{}, &Order{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}
