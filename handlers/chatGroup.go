package handlers

import (
	"net/http"

	"chat-backend/models"
	"chat-backend/services"

	"github.com/gin-gonic/gin"
)

// CreateChatHandler handles creating a new chat group
func CreateChatHandler(c *gin.Context) {
	var chat models.ChatGroup

	// Bind JSON data
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create chat
	err := services.CreateChatGroup(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully"})
}
