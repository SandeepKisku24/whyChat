package services

import (
	"chat-backend/config"
	"chat-backend/models"
	"log"

	"github.com/google/uuid"
)

// for a global group chat
const globalChatID = "global_chat" // Single group chat ID

// CreateUser adds a user to Firestore and ensures they join the group chat
func CreateUser(phoneNumber string, name string) (*models.User, error) {
	client, ctx, err := config.GetFirestoreClient()
	if err != nil {
		return nil, err
	}

	// Check if the user already exists
	userRef := client.Collection("users").Doc(phoneNumber)
	doc, err := userRef.Get(ctx)
	if err == nil && doc.Exists() {
		// User already exists
		var existingUser models.User
		doc.DataTo(&existingUser)
		return &existingUser, nil
	}

	// New user
	user := models.User{
		UserID:      phoneNumber,
		Name:        name,
		PhoneNumber: phoneNumber,
		ChatGroups:  []string{globalChatID}, // Auto-add to global chat
	}

	// Save to Firestore
	_, err = userRef.Set(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	log.Println("User created successfully:", phoneNumber)
	return &user, nil
}

// CreateUser adds a new user to Firestore for a new chat Group
func CreateNewUser(user models.User) error {
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
