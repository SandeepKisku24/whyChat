package services

import (
	"chat-backend/config"
	"chat-backend/models"
	"log"
	"time"

	"github.com/google/uuid"
)

// SendMessage saves a new message to Firestore
func SendMessage(senderID, chatGroupID, messageText string) (*models.Message, error) {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		return nil, err
	}

	// Create new message
	messageID := uuid.New().String()
	message := models.Message{
		MessageID:   messageID,
		ChatGroupID: chatGroupID,
		SenderID:    senderID,
		Message:     messageText,
		Timestamp:   time.Now(),
		MessageType: "text",
	}

	// Save message to Firestore
	_, err = client.Collection("messages").Doc(messageID).Set(ctx, message)
	if err != nil {
		log.Printf(" Error sending message: %v", err)
		return nil, err
	}

	log.Printf("Message sent successfully to chat group [%s]: %s", chatGroupID, messageText)
	return &message, nil
}

// GetMessages retrieves all messages from the global chat, sorted by timestamp
func GetMessages(chatGroupID string) ([]models.Message, error) {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		return nil, err
	}

	log.Printf("🔍 Fetching messages for chat group: %s", chatGroupID)

	// Query messages in the specified chat group, ordered by timestamp
	query := client.Collection("messages").Where("chatGroupId", "==", chatGroupID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Firestore query error: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d messages from chat group [%s]", len(docs), chatGroupID)

	// Parse messages
	var messages []models.Message
	for _, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}
