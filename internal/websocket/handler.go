package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now (can be restricted in production)
		return true
	},
}

// Handler creates a WebSocket handler for the given hub
func Handler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get client ID and user ID from query params
		clientID := r.URL.Query().Get("clientId")
		if clientID == "" {
			clientID = "anonymous_" + r.RemoteAddr
		}
		userID := r.URL.Query().Get("userId")

		// Upgrade connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		// Create new client
		client := NewClient(hub, conn, clientID, userID)

		// Register client with hub
		hub.register <- client

		// Start goroutines for reading and writing
		go client.WritePump()
		go client.ReadPump()

		log.Printf("WebSocket client connected: %s (user: %s)", clientID, userID)
	}
}
