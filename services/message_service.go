package services

import (
	"chat-backend/config"
	"chat-backend/models"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// SendMessage saves the message in Firestore
func SendMessage(senderID, chatGroupID, message string) (*models.Message, error) {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		log.Println("Firestore client error:", err)
		return nil, err
	}

	// Generate unique message ID
	messageID := uuid.New().String()

	// Create message object
	msg := models.Message{
		MessageID:   messageID,
		ChatGroupID: chatGroupID,
		SenderID:    senderID,
		Message:     message,
		Timestamp:   time.Now(),
		MessageType: "text",
	}

	// Store message in Firestore
	_, err = client.Collection("messages").Doc(messageID).Set(ctx, msg)
	if err != nil {
		log.Println("Failed to store message:", err)
		return nil, err
	}

	// Reference to chat group document
	chatGroupRef := client.Collection("chat_groups").Doc(chatGroupID)

	// Check if chat group exists
	_, err = chatGroupRef.Get(ctx)
	if err != nil {
		// If the document doesn't exist, create it with an empty MessageIDs array
		_, err = chatGroupRef.Set(ctx, map[string]interface{}{
			"MessageIDs": []string{messageID}, // Initialize with the first message
		})
		if err != nil {
			log.Println("Failed to create chat group:", err)
			return nil, err
		}
	} else {
		// If chat group exists, update it normally
		_, err = chatGroupRef.Update(ctx, []firestore.Update{
			{Path: "MessageIDs", Value: firestore.ArrayUnion(messageID)},
		})
		if err != nil {
			log.Println("Failed to update chat group:", err)
			return nil, err
		}

	}

	return &msg, nil
}

// GetMessages retrieves messages from Firestore
func GetMessages(chatGroupID string) ([]models.Message, error) {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		log.Println("Firestore client error:", err)
		return nil, err
	}

	// Query messages
	iter := client.Collection("messages").
		Where("chatGroupId", "==", chatGroupID). // Ensure exact match
		OrderBy("timestamp", firestore.Asc).
		Documents(ctx)

	var messages []models.Message
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Failed to iterate messages:", err)
			return nil, err
		}

		var msg models.Message
		if err := doc.DataTo(&msg); err != nil {
			log.Println("Failed to parse message data:", err)
			continue
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func DeleteMessage(messageID string) error {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		log.Println(" FireStore client error")
		return err
	}

	msgRef := client.Collection("messages").Doc(messageID)

	_, err = msgRef.Update(ctx, []firestore.Update{
		{Path: "isDeleted", Value: true},
	})

	if err != nil {
		log.Println("the message is not deleted")
		return err
	}

	return nil

}
