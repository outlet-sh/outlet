package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB
)

// Client represents a WebSocket client connection
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Client ID
	ID string

	// User ID (if authenticated)
	UserID string

	// Organization IDs this client is subscribed to
	orgSubscriptions map[string]bool
	orgMu            sync.RWMutex
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, clientID, userID string) *Client {
	return &Client{
		hub:              hub,
		conn:             conn,
		send:             make(chan []byte, 256),
		ID:               clientID,
		UserID:           userID,
		orgSubscriptions: make(map[string]bool),
	}
}

// SubscribeToOrg subscribes the client to updates for an organization
func (c *Client) SubscribeToOrg(orgID string) {
	c.orgMu.Lock()
	defer c.orgMu.Unlock()
	c.orgSubscriptions[orgID] = true
}

// UnsubscribeFromOrg unsubscribes the client from updates for an organization
func (c *Client) UnsubscribeFromOrg(orgID string) {
	c.orgMu.Lock()
	defer c.orgMu.Unlock()
	delete(c.orgSubscriptions, orgID)
}

// IsSubscribedToOrg checks if the client is subscribed to an organization
func (c *Client) IsSubscribedToOrg(orgID string) bool {
	c.orgMu.RLock()
	defer c.orgMu.RUnlock()
	return c.orgSubscriptions[orgID]
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse the message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to parse WebSocket message: %v", err)
			continue
		}

		// Handle message types
		switch msg.Type {
		case TypePing:
			// Respond with pong
			pongMsg := NewMessage(TypePong, nil)
			if data, err := json.Marshal(pongMsg); err == nil {
				c.send <- data
			}

		case TypeSubscribe:
			// Subscribe to org updates
			if orgID, ok := msg.Data.(string); ok {
				c.SubscribeToOrg(orgID)
				log.Printf("Client %s subscribed to org %s", c.ID, orgID)
			} else if dataMap, ok := msg.Data.(map[string]interface{}); ok {
				if orgID, ok := dataMap["org_id"].(string); ok {
					c.SubscribeToOrg(orgID)
					log.Printf("Client %s subscribed to org %s", c.ID, orgID)
				}
			}

		case TypeUnsubscribe:
			// Unsubscribe from org updates
			if orgID, ok := msg.Data.(string); ok {
				c.UnsubscribeFromOrg(orgID)
				log.Printf("Client %s unsubscribed from org %s", c.ID, orgID)
			} else if dataMap, ok := msg.Data.(map[string]interface{}); ok {
				if orgID, ok := dataMap["org_id"].(string); ok {
					c.UnsubscribeFromOrg(orgID)
					log.Printf("Client %s unsubscribed from org %s", c.ID, orgID)
				}
			}
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
