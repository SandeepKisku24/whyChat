package models

type Message struct {
	MessageID   string `firestore:"messageId"`
	ChatGroupID string `firestore:"chatGroupId"`
	SenderID    string `firestore:"senderId"`
	Message     string `firestore:"message"`
	Timestamp   int64  `firestore:"timestamp"`
	MessageType string `firestore:"messageType"`
}
