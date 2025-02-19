package handlers

import (
	"net/http"

	"chat-backend/models"
	"chat-backend/services"

	"github.com/gin-gonic/gin"
)

// for a global chat group

// CreateUserHandler creates a new user
func CreateUserHandler(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Name        string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := services.CreateUser(request.PhoneNumber, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "user": user})
}

// *********************

// CreateUserHandler handles user creation for new chatGroup
func CreateNewUserHandler(c *gin.Context) {
	var user models.User

	// Bind JSON data
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create user
	err := services.CreateNewUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
