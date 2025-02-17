package models

type User struct {
	UserID      string   `firestore:"user_id"`
	Name        string   `firestore:"name"`
	PhoneNumber string   `firestore:"phoneNumber"`
	ChatGroups  []string `firestore:"chatGroups"`
}
