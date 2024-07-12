package controllers

import (
	"e-commerce/db"
	"e-commerce/dto"
	"e-commerce/models"
	"e-commerce/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user
// @Summary Register a new user
// @Description Create a new user with the given details
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.UserRegisterRequest true "User information"
// @Success 201 {object} dto.UserRegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/register [post]
func RegisterUser(c *gin.Context) {
	var userInput dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to hash password"})
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
	}

	// Prepare response
	userResponse := dto.UserRegisterResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  newUser.Role,
	}

	c.JSON(http.StatusCreated, userResponse)
}

// LoginUser handles user login
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.UserLoginRequest true "User credentials"
// @Success 200 {object} dto.UserLoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/login [post]
func LoginUser(c *gin.Context) {
	var userInput dto.UserLoginRequest
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Retrieve user by email
	var existingUser models.User
	if err := db.DB.Where("email = ?", userInput.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Compare stored hashed password with input password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Generate JWT token with user role
	token, err := utils.GenerateJWT(existingUser.ID, existingUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Unable to generate token"})
		return
	}

	// Prepare response
	userResponse := dto.UserRegisterResponse{
		ID:    existingUser.ID,
		Name:  existingUser.Name,
		Email: existingUser.Email,
		Role:  existingUser.Role,
	}

	loginResponse := dto.UserLoginResponse{
		User:  userResponse,
		Token: token,
	}

	c.JSON(http.StatusOK, loginResponse)
}
