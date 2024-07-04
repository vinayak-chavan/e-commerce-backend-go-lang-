package controllers

import (
	"e-commerce/db"
	"e-commerce/models"
	"e-commerce/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user
func RegisterUser(c *gin.Context) {
	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user record
	newUser := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: string(hashedPassword),
		Role:     userInput.Role,
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Remove password from response for security reasons
	newUser.Password = ""

	c.JSON(http.StatusCreated, newUser)
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user by email
	var existingUser models.User
	if err := db.DB.Where("email = ?", userInput.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare stored hashed password with input password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token with user role
	token, err := utils.GenerateJWT(existingUser.ID, existingUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token"})
		return
	}

	// Remove password from response for security reasons
	existingUser.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user":  existingUser,
		"token": token,
	})
}
