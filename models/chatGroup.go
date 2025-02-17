package models

type ChatGroup struct {
	ChatGroupID string   `firestore:"chatGroupId"`
	Members     []string `firestore:"members"`
	MessageIDs  []string `firestore:"messageIds"`
}
