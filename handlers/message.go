package handlers

import (
	"net/http"

	"chat-backend/models"
	"chat-backend/services"

	"github.com/gin-gonic/gin"
)

// SendMessageHandler handles sending a new message
func SendMessageHandler(c *gin.Context) {
	var message models.Message

	// Bind JSON data
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to send message
	err := services.SendMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
