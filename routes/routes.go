package routes

import (
	"chat-backend/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.Engine) {
	router.POST("/users", handlers.CreateUserHandler)
	router.POST("/chats", handlers.CreateChatHandler)
	router.POST("/messages", handlers.SendMessageHandler)

	router.GET("/messages", handlers.GetMessagesHandler)
	// OTP Routes
	// router.POST("/generate-otp", handlers.GenerateOTP)
	// router.POST("/verify-otp", handlers.VerifyOTP)
}
