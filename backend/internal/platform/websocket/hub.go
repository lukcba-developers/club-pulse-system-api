package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
)

// Config holds WebSocket server configuration.
type Config struct {
	WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int64
}

// DefaultConfig provides standard operational values.
var DefaultConfig = Config{
	WriteWait:      10 * time.Second,
	PongWait:       60 * time.Second,
	PingPeriod:     54 * time.Second, // Must be less than PongWait
	MaxMessageSize: 512,
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Security: In production, check r.Header.Get("Origin") against allowed domains.
		return true
	},
}

// Client represents a connected WebSocket user.
// Refactored to be cleaner (Logic moved to Hub).
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

// Hub maintains the set of active clients and handles message routing.
// Architected for High Performance with O(1) user lookups and Topic Subscriptions.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// O(1) lookup for User -> Clients (User can have multiple devices).
	userClients map[string]map[*Client]bool

	// Topic subscriptions for filtered broadcasting.
	// topic -> client -> true
	subscriptions map[string]map[*Client]bool

	// Inbound messages from clients (Requests).
	commands chan clientCommand

	// Outbound messages to broadcast/multicast.
	broadcast chan broadcastMessage

	// Targeted messages (Unicast).
	unicast chan unicastMessage

	// Register/Unregister requests.
	register   chan *Client
	unregister chan *Client

	// Sync primitives (Though we try to run single-threaded loop).
	// Locking strategy: The Hub.Run loop owns all maps.
	// External access (like SendToUser called from HTTP handler) needs thread-safe injection.
	// We use channels for everything to stay lock-free in the main loop.
}

type clientCommand struct {
	client  *Client
	payload []byte
}

type broadcastMessage struct {
	topic   string // If empty, global broadcast.
	payload []byte
	exclude *Client
}

// NewHub initializes the Hub with optimized data structures.
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		userClients:   make(map[string]map[*Client]bool),
		subscriptions: make(map[string]map[*Client]bool),
		broadcast:     make(chan broadcastMessage, 256),
		unicast:       make(chan unicastMessage, 256),
		commands:      make(chan clientCommand, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}
}

// Run handles all Hub events in a single goroutine to avoid race conditions and mutex contention.
func (h *Hub) Run(ctx context.Context) {
	// Background Redis Subscription
	go h.subscribeToRedis(ctx)

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case cmd := <-h.commands:
			h.processClientCommand(cmd.client, cmd.payload)

		case msg := <-h.broadcast:
			h.dispatchMessage(msg)

		case msg := <-h.unicast:
			h.dispatchUnicast(msg)

		case <-ctx.Done():
			h.shutdown()
			return
		}
	}
}

// --- Internal Logic (Single Threaded - Safe) ---

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true

	// Add to User Map
	if _, ok := h.userClients[client.userID]; !ok {
		h.userClients[client.userID] = make(map[*Client]bool)
	}
	h.userClients[client.userID][client] = true

	log.Printf("Client connected: %s", client.userID)
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		// Clean up global list
		delete(h.clients, client)

		// Clean up user map
		if _, ok := h.userClients[client.userID]; ok {
			delete(h.userClients[client.userID], client)
			if len(h.userClients[client.userID]) == 0 {
				delete(h.userClients, client.userID)
			}
		}

		// Clean up subscriptions
		// Performance Note: Iterate topics is O(T). If T is large, invert this map (Client -> Topics).
		// For MVP/Medium scale, this is acceptable.
		for topic, subs := range h.subscriptions {
			if _, ok := subs[client]; ok {
				delete(subs, client)
				if len(subs) == 0 {
					delete(h.subscriptions, topic)
				}
			}
		}

		close(client.send)
	}
}

func (h *Hub) processClientCommand(client *Client, payload []byte) {
	var msg struct {
		Action  string   `json:"action"`
		Targets []string `json:"targets"`
	}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return // Ignore malformed
	}

	switch msg.Action {
	case "subscribe":
		for _, topic := range msg.Targets {
			if h.subscriptions[topic] == nil {
				h.subscriptions[topic] = make(map[*Client]bool)
			}
			h.subscriptions[topic][client] = true
		}
	case "unsubscribe":
		for _, topic := range msg.Targets {
			if subs, ok := h.subscriptions[topic]; ok {
				delete(subs, client)
			}
		}
	}
}

func (h *Hub) dispatchMessage(msg broadcastMessage) {
	if msg.topic != "" {
		// Topic Broadcast
		if subs, ok := h.subscriptions[msg.topic]; ok {
			for client := range subs {
				if client != msg.exclude {
					h.sendSafe(client, msg.payload)
				}
			}
		}
	} else {
		// Global Broadcast
		for client := range h.clients {
			h.sendSafe(client, msg.payload)
		}
	}
}

func (h *Hub) dispatchUnicast(msg unicastMessage) {
	if clients, ok := h.userClients[msg.targetUserID]; ok {
		for client := range clients {
			h.sendSafe(client, msg.payload)
		}
	}
}

func (h *Hub) sendSafe(client *Client, payload []byte) {
	select {
	case client.send <- payload:
	default:
		// Slow client, drop or kick. For now, we drop message, kicking is handled by heartbeats ideally.
		// h.unregisterClient(client) // Auto-kick policy could go here
	}
}

func (h *Hub) shutdown() {
	for client := range h.clients {
		h.unregisterClient(client)
	}
}

// --- External API (Thread-Safe via Channels) ---

// SendToUser sends a message efficiently to a specific user across all their devices.
// This command is technically external, but we inject it into the main loop via a special "Internal" command
// or we just handle it here if we want to block? No, blocking is bad.
// Ideally, `broadcast` needs to support "TargetUser".
// For now, let's just cheat and assume direct access if we can't change the Hub struct too much,
// BUT since we are Architects, let's fix strictly.
// We'll add a 'DirectMessage' channel or use Broadcast with metadata.
// For simplicity in this Refactor Step 1, I'll add a 'unicast' channel/method logic.
// However, since `SendToUser` is exported, let's make it push to a channel.
// We need to extend the `Hub` to support this properly.
// Let's add `unicast chan unicastMessage` to Hub struct? Or reuse `broadcast` with target.
// Simpler: Just generic `events` channel.
// To keep diff minimal, I will implement a safe closure.
func (h *Hub) SendToUser(userID string, message []byte) {
	// We run a goroutine to not block the caller, pushing to a channel processed by Run loop.
	// But `Hub` struct definition needs to change to hold this channel.
	// I'll skip adding a new channel to keep struct simple and assume this method is infrequently called or change design to use `broadcast`.

	// Locking strategy revisited: If we want O(1) read, we need RLock.
	// But `Run` uses no mutex. We cannot mix strategies safely without great care.
	// Correct approach: `Run` owns the map. `SendToUser` MUST send a message to `Run`.
	// Using a `inject` channel.

	// Since I can't effectively add a new channel and restart the `Run` loop in a live/hot reload sense easily without restarting app,
	// I will stick to the existing `broadcast` channel pattern but perhaps abuse it?
	// No, clean code. I will check userClients in the loop.

	// WAIT: I added `userClients` to `Hub`. I can't read it here safely while `Run` writes it.
	// I MUST assume `SendToUser` is not thread safe unless I add a Mutex OR use channels.
	// Given strict requirements: I will stick to the "Run Loop Owns All" pattern.
	// I will add `unicast` channel to Hub.

	// Note: Since I am rewriting the file, I CAN change the struct.
	// I will add `unicast chan unicastMessage`.

	h.unicast <- unicastMessage{targetUserID: userID, payload: message}
}

// Needs to be added to struct and constructor above.
// I will edit the `Hub` struct in the big file write.

type unicastMessage struct {
	targetUserID string
	payload      []byte
}

// --- Redis Integration ---

func (h *Hub) subscribeToRedis(ctx context.Context) {
	subscriber := platformRedis.NewEventSubscriber()

	// This fits the "Adapter" pattern logic.
	err := subscriber.SubscribeToBookings(ctx, func(event platformRedis.BookingEvent) {
		msg := WebSocketMessage{
			Type:      event.Type,
			Payload:   event,
			Timestamp: time.Now(),
		}
		data, _ := json.Marshal(msg)

		// Determine topic based on event? e.g. "facility:{id}"
		// For now, broadcast global or specific?
		// Assuming global for simplified MVP logic as per original file.
		// Or smart routing:
		// h.broadcast <- broadcastMessage{topic: "facility:" + event.FacilityID, payload: data}

		// Maintaining original "Broadcast All" behavior for compatibility,
		// but enabling future smart routing.
		h.broadcast <- broadcastMessage{payload: data}
	})

	if err != nil {
		log.Printf("Redis error: %v", err)
	}
}

// --- Handlers ---

// HandleWebSocket upgrades HTTP to WS.
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			userID = "anonymous"
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256), // Buffered to handle bursts
			userID: userID,
		}

		client.hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}

// writePump writes messages to the websocket.
func (c *Client) writePump() {
	ticker := time.NewTicker(DefaultConfig.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(DefaultConfig.WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Optimized: Flush queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(DefaultConfig.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump reads messages from the websocket.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(DefaultConfig.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(DefaultConfig.PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(DefaultConfig.PongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS Error: %v", err)
			}
			break
		}
		// Send command to Hub (Thread-safe state mutation)
		c.hub.commands <- clientCommand{client: c, payload: message}
	}
}

// --- DTOs ---

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}
