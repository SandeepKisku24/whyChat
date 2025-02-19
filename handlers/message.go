package handlers

import (
	"chat-backend/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

// SendMessageHandler allows users to send messages to the group
func SendMessageHandler(c *gin.Context) {
	var request struct {
		SenderID    string `json:"sender_id" binding:"required"`
		ChatGroupID string `json:"chat_group_id" binding:"required"`
		Message     string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	message, err := services.SendMessage(request.SenderID, request.ChatGroupID, request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent", "data": message})
}

// GetMessagesHandler retrieves all messages from the global chat
func GetMessagesHandler(c *gin.Context) {
	chatGroupID := c.Query("chat_group_id")
	if chatGroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat group ID is required"})
		return
	}

	messages, err := services.GetMessages(chatGroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
