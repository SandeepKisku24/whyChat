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

type MessageResponse struct {
	MessageID   string    `json:"messageId"`
	ChatGroupID string    `json:"chatGroupId"`
	SenderID    string    `json:"senderId"`
	SenderName  string    `json:"senderName"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"messageType"`
}
