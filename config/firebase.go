package config

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var client *firestore.Client
var ctx context.Context

// InitFirestore initializes Firestore client
func InitFirestore() error {
	ctx = context.Background()

	// ✅ Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// ✅ Get environment variables
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS")

	// ✅ Validate environment variables
	if projectID == "" || credentialsPath == "" {
		log.Fatal("Missing FIREBASE_PROJECT_ID or FIREBASE_CREDENTIALS in environment variables")
		return errors.New("missing Firebase config")
	}

	// ✅ Convert relative path to absolute
	absCredentialsPath, err := filepath.Abs(credentialsPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
		return err
	}

	// ✅ Initialize Firestore
	sa := option.WithCredentialsFile(absCredentialsPath)
	client, err = firestore.NewClient(ctx, projectID, sa)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return err
	}

	log.Println("Firestore initialized successfully!")
	return nil
}

// GetFirestoreClient safely retrieves the Firestore client
func GetFirestoreClient() (*firestore.Client, context.Context, error) {
	if client == nil {
		return nil, nil, errors.New("Firestore client is not initialized")
	}
	return client, ctx, nil
}

// CloseFirestoreClient closes Firestore client connection
func CloseFirestoreClient() {
	if client != nil {
		client.Close()
		log.Println("Firestore client closed.")
	}
}
