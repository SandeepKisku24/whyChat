package models

import "time"

type Message struct {
	MessageID   string    `firestore:"messageId"`
	ChatGroupID string    `firestore:"chatGroupId"`
	SenderID    string    `firestore:"senderId"`
	Message     string    `firestore:"message"`
	Timestamp   time.Time `json:"timestamp" firestore:"timestamp"`
	MessageType string    `firestore:"messageType"`
}
