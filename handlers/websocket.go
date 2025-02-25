package handlers

import (
	"chat-backend/services"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Store connections for each chat room
var chatRooms = make(map[string]map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

// Handle WebSocket connections
func HandleWebSocket(c *gin.Context) {
	chatGroupID := c.Query("chat_group_id")
	if chatGroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat group ID is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer removeConnection(chatGroupID, conn) // Ensure cleanup on exit

	// Set initial read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Handle Pong response
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Extend timeout on Pong
		return nil
	})

	// Register connection
	mutex.Lock()
	if chatRooms[chatGroupID] == nil {
		chatRooms[chatGroupID] = make(map[*websocket.Conn]bool)
	}
	chatRooms[chatGroupID][conn] = true
	mutex.Unlock()

	log.Println("New WebSocket connection for chat group:", chatGroupID)

	// Start pinging the client every 30 seconds
	go pingClient(chatGroupID, conn)

	// Read messages loop
	for {
		var msg struct {
			SenderID string `json:"SenderID"`
			Message  string `json:"Message"`
		}

		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket Read Error:", err)
			}
			break // Exit loop and trigger cleanup
		}

		// Refresh deadline on every received message
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Store and broadcast message
		storedMessage, err := services.SendMessage(msg.SenderID, chatGroupID, msg.Message)
		if err != nil {
			log.Println("Failed to store message:", err)
			continue
		}

		log.Println("Message stored successfully:", storedMessage) // Log when message is stored

		broadcastMessage(chatGroupID, storedMessage)
	}
}

// Ping client every 30 seconds to keep connection alive
func pingClient(chatGroupID string, conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mutex.Lock()
		if _, exists := chatRooms[chatGroupID][conn]; !exists {
			mutex.Unlock()
			return // Stop pinging if connection is removed
		}
		mutex.Unlock()

		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Println("Ping failed, closing connection:", err)
			removeConnection(chatGroupID, conn)
			return // Stop pinging on failure
		}
	}
}

// Remove disconnected WebSocket connections
func removeConnection(chatGroupID string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	// Remove disconnected connection
	if _, exists := chatRooms[chatGroupID]; exists {
		if _, connExists := chatRooms[chatGroupID][conn]; connExists {
			delete(chatRooms[chatGroupID], conn)
			conn.Close() // Ensure connection is properly closed
			log.Println("User disconnected from chat group:", chatGroupID)
		}
	}

	// Only delete chat group if it's truly empty
	if len(chatRooms[chatGroupID]) == 0 {
		delete(chatRooms, chatGroupID)
		log.Println("Closed chat group:", chatGroupID)
	}
}

// Broadcast message to all clients
func broadcastMessage(chatGroupID string, message interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	connections, exists := chatRooms[chatGroupID]
	if !exists {
		return
	}

	log.Println("Broadcasting message to chat group:", chatGroupID, "Message:", message) // Log before broadcasting

	for conn := range connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Println("Error broadcasting message:", err)
			removeConnection(chatGroupID, conn) // Properly remove failed connections
		}
	}
}
