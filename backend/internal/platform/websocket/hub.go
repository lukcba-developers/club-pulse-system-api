package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, validate origin properly
		return true
	},
}

// Client represents a WebSocket client connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   string
	channels []string // Subscribed channels (e.g., "facility:uuid")
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex

	// Redis subscriber for distributed messaging
	redis *platformRedis.RedisClient
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      platformRedis.GetClient(),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run(ctx context.Context) {
	// Start Redis subscription in background
	go h.subscribeToRedis(ctx)

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected: %s", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected: %s", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()

		case <-ctx.Done():
			return
		}
	}
}

// subscribeToRedis listens for Redis Pub/Sub events and broadcasts to clients
func (h *Hub) subscribeToRedis(ctx context.Context) {
	subscriber := platformRedis.NewEventSubscriber()

	// Subscribe to booking events
	err := subscriber.SubscribeToBookings(ctx, func(event platformRedis.BookingEvent) {
		// Convert to WebSocket message
		msg := WebSocketMessage{
			Type:      event.Type,
			Payload:   event,
			Timestamp: time.Now(),
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal WebSocket message: %v", err)
			return
		}

		// Broadcast to all connected clients
		h.broadcast <- data
	})

	if err != nil {
		log.Printf("Failed to subscribe to Redis: %v", err)
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.userID == userID {
			select {
			case client.send <- message:
			default:
				// Buffer full, skip
			}
		}
	}
}

// WebSocketMessage represents a message sent to clients
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// HandleWebSocket handles WebSocket upgrade and connection
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from auth context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			userID = "anonymous"
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			userID: userID.(string),
		}

		hub.register <- client

		// Start goroutines for reading and writing
		go client.writePump()
		go client.readPump()
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
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
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
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

		// Handle client messages (e.g., subscribe to specific facilities)
		c.handleMessage(message)
	}
}

// handleMessage handles incoming messages from clients
func (c *Client) handleMessage(message []byte) {
	var msg struct {
		Action  string   `json:"action"`
		Targets []string `json:"targets"`
	}

	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}

	switch msg.Action {
	case "subscribe":
		c.channels = append(c.channels, msg.Targets...)
	case "unsubscribe":
		// Remove targets from channels
		newChannels := make([]string, 0)
		for _, ch := range c.channels {
			found := false
			for _, t := range msg.Targets {
				if ch == t {
					found = true
					break
				}
			}
			if !found {
				newChannels = append(newChannels, ch)
			}
		}
		c.channels = newChannels
	}
}
