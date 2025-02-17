package main

import (
	"chat-backend/config"
	"chat-backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting server...")

	// ✅ Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// ✅ Check if required environment variables are set
	if os.Getenv("FIREBASE_PROJECT_ID") == "" || os.Getenv("FIREBASE_CREDENTIALS") == "" {
		log.Fatal("Missing FIREBASE_PROJECT_ID or FIREBASE_CREDENTIALS in environment variables")
	}

	// ✅ Initialize Firestore before using it
	if err := config.InitFirestore(); err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}

	// Setup router
	r := gin.Default()
	routes.RegisterRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
