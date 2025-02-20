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

// GetMessages retrieves all messages and includes the sender's name
func GetMessages(chatGroupID string) ([]models.MessageResponse, error) {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		return nil, err
	}

	log.Printf("🔍 Fetching messages for chat group: %s", chatGroupID)

	// Query messages in the specified chat group
	query := client.Collection("messages").Where("chatGroupId", "==", chatGroupID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Firestore query error: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d messages from chat group [%s]", len(docs), chatGroupID)

	// Parse messages and fetch sender names
	var messages []models.MessageResponse
	for _, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Fetch sender's name from the users collection
		userDoc, err := client.Collection("users").Doc(message.SenderID).Get(ctx)
		var senderName string
		if err != nil {
			log.Printf("Error fetching sender name for ID %s: %v", message.SenderID, err)
			senderName = "Unknown" // Fallback if user not found
		} else {
			var user models.User
			if err := userDoc.DataTo(&user); err == nil {
				senderName = user.Name
			} else {
				senderName = "Unknown"
			}
		}

		// Append message with sender name
		messages = append(messages, models.MessageResponse{
			MessageID:   message.MessageID,
			ChatGroupID: message.ChatGroupID,
			SenderID:    message.SenderID,
			SenderName:  senderName,
			Message:     message.Message,
			Timestamp:   message.Timestamp,
			MessageType: message.MessageType,
		})
	}

	return messages, nil
}
