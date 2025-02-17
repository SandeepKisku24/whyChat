package services

import (
	"chat-backend/config"
	"chat-backend/models"
	"log"

	"github.com/google/uuid"
)

// CreateChatGroup creates a new chat group
func CreateChatGroup(chat models.ChatGroup) error {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		log.Printf("Firestore client error: %v", err)
		return err
	}

	// 🔹 Generate a unique ChatGroupID if not provided
	if chat.ChatGroupID == "" {
		chat.ChatGroupID = uuid.New().String()
	}

	_, err = client.Collection("chats").Doc(chat.ChatGroupID).Set(ctx, chat)
	if err != nil {
		log.Printf("Error creating chat group: %v", err)
		return err
	}

	log.Println("Chat group created successfully with ID:", chat.ChatGroupID)
	return nil
}
