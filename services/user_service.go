package services

import (
	"chat-backend/config"
	"chat-backend/models"
	"log"

	"github.com/google/uuid"
)

// CreateUser adds a new user to Firestore
func CreateUser(user models.User) error {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		log.Printf("Firestore client error: %v", err)
		return err
	}

	// 🔹 Generate a unique ID if user.UserID is empty
	if user.UserID == "" {
		user.UserID = uuid.New().String()
	}

	_, err = client.Collection("users").Doc(user.UserID).Set(ctx, user)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		return err
	}

	log.Println("User added successfully with ID:", user.UserID)
	return nil
}
