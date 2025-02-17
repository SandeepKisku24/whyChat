package services

import (
	"chat-backend/config"
	"chat-backend/models"
	"log"

	"github.com/google/uuid"
)

// SendMessage stores a message in Firestore
func SendMessage(message models.Message) error {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		log.Printf("Firestore client error: %v", err)
		return err
	}

	// 🔹 Generate a unique MessageID if not provided
	if message.MessageID == "" {
		message.MessageID = uuid.New().String()
	}

	_, err = client.Collection("messages").Doc(message.MessageID).Set(ctx, message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	log.Println("Message sent successfully with ID:", message.MessageID)
	return nil
}
