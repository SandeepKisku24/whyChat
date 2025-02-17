package handlers

import (
	"net/http"

	"chat-backend/models"
	"chat-backend/services"

	"github.com/gin-gonic/gin"
)

// CreateUserHandler handles user creation
func CreateUserHandler(c *gin.Context) {
	var user models.User

	// Bind JSON data
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create user
	err := services.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
