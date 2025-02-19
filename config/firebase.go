package config

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var client *firestore.Client
var ctx context.Context
var firebaseApp *firebase.App
var authClient *auth.Client

// InitFirestore initializes Firestore and Firebase Authentication clients
func InitFirestore() error {
	ctx = context.Background()

	// Get the service account credentials path and Firebase API key from environment variables
	serviceAccountPath := os.Getenv("FIREBASE_CREDENTIALS")
	if serviceAccountPath == "" {
		log.Fatal("FIREBASE_CREDENTIALS environment variable is not set")
		return errors.New("FIREBASE_CREDENTIALS not set")
	}

	// Initialize Firebase app with the credentials file
	opt := option.WithCredentialsFile(serviceAccountPath)
	var err error
	firebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
		return err
	}

	// Initialize Firestore client
	client, err = firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return err
	}

	// Initialize Firebase Authentication client
	authClient, err = firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firebase Auth client: %v", err)
		return err
	}

	log.Println("Firebase and Firestore initialized successfully!")
	return nil
}

// GetFirestoreClient safely retrieves the Firestore client
func GetFirestoreClient() (*firestore.Client, context.Context, error) {
	if client == nil {
		return nil, nil, errors.New("Firestore client is not initialized")
	}
	return client, ctx, nil
}

// GetAuthClient safely retrieves the Firebase Auth client
func GetAuthClient() *auth.Client {
	if authClient == nil {
		log.Fatalf("Firebase Auth client is not initialized")
	}
	return authClient
}

// CloseFirestoreClient closes Firestore client connection
func CloseFirestoreClient() {
	if client != nil {
		client.Close()
		log.Println("Firestore client closed.")
	}
}
