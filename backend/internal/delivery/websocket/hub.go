package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yourusername/bus-booking/internal/entities"
	"github.com/yourusername/bus-booking/internal/repositories/cache"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Hub maintains active WebSocket connections and broadcasts messages
type Hub struct {
	// TripID -> Set of client connections
	rooms      map[string]map[*Client]bool
	roomsMutex sync.RWMutex

	// Redis cache for pub/sub
	cache *cache.RedisCache
}

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	tripID string
}

// NewHub creates a new WebSocket hub
func NewHub(cache *cache.RedisCache) *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]bool),
		cache: cache,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	tripID := r.URL.Query().Get("trip_id")
	if tripID == "" {
		conn.Close()
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		tripID: tripID,
	}

	h.register(client)
	defer h.unregister(client)

	// Start goroutines for reading and writing
	go client.writePump()
	client.readPump()
}

func (h *Hub) register(client *Client) {
	h.roomsMutex.Lock()
	defer h.roomsMutex.Unlock()

	if h.rooms[client.tripID] == nil {
		h.rooms[client.tripID] = make(map[*Client]bool)

		// Subscribe to Redis pub/sub for this trip
		go h.subscribeToTripUpdates(client.tripID)
	}
	h.rooms[client.tripID][client] = true
}

func (h *Hub) unregister(client *Client) {
	h.roomsMutex.Lock()
	defer h.roomsMutex.Unlock()

	if clients, ok := h.rooms[client.tripID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)

			// Clean up empty rooms
			if len(clients) == 0 {
				delete(h.rooms, client.tripID)
			}
		}
	}
}

func (h *Hub) broadcast(tripID string, message []byte) {
	h.roomsMutex.RLock()
	defer h.roomsMutex.RUnlock()

	if clients, ok := h.rooms[tripID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

func (h *Hub) subscribeToTripUpdates(tripID string) {
	ctx := context.Background()
	tripUUID, err := uuid.Parse(tripID)
	if err != nil {
		return
	}

	pubsub := h.cache.SubscribeTripUpdates(ctx, tripUUID)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		// Broadcast to all clients in this room
		var seat entities.SeatInfo
		if err := json.Unmarshal([]byte(msg.Payload), &seat); err == nil {
			message, _ := json.Marshal(map[string]interface{}{
				"type": "seat_update",
				"data": seat,
			})
			h.broadcast(tripID, message)
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Handle incoming messages (e.g., heartbeat)
		log.Printf("Received message: %s", message)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

// Global hub instance
var globalHub *Hub

// StartWebSocketServer initializes the WebSocket server
func StartWebSocketServer(cache *cache.RedisCache) {
	globalHub = NewHub(cache)

	http.HandleFunc("/ws", globalHub.HandleWebSocket)

	log.Println("WebSocket server starting on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("WebSocket server error:", err)
	}
}
