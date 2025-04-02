package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Client represents a connected WebSocket client
type Client struct {
	Conn     *websocket.Conn
	Username string
}

// ClientManager manages WebSocket clients and broadcasts
type ClientManager struct {
	clients    map[*websocket.Conn]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

// Message represents a chat message
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Sender  string `json:"sender"`
	Time    string `json:"time"`
}

// NewClientManager creates a new client manager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*websocket.Conn]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *websocket.Conn),
	}
}

// Start begins the client manager processing loop
func (m *ClientManager) Start() {
	go func() {
		for {
			select {
			case client := <-m.register:
				m.mutex.Lock()
				m.clients[client.Conn] = client
				m.mutex.Unlock()
				log.Printf("Client connected: %s", client.Username)

				// Announce new user to all clients
				joinMsg := fmt.Sprintf("System: %s joined the chat", client.Username)
				m.broadcast <- []byte(joinMsg)

			case conn := <-m.unregister:
				m.mutex.Lock()
				if client, ok := m.clients[conn]; ok {
					username := client.Username
					delete(m.clients, conn)
					conn.Close()

					// Announce user departure
					leaveMsg := fmt.Sprintf("System: %s left the chat", username)
					m.broadcast <- []byte(leaveMsg)
				}
				m.mutex.Unlock()
				log.Println("Client disconnected")

			case message := <-m.broadcast:
				m.mutex.Lock()
				for conn, _ := range m.clients {
					if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
						log.Println("Write error:", err)
						conn.Close()
						delete(m.clients, conn)
					}
				}
				m.mutex.Unlock()
			}
		}
	}()
}

// Run starts the Fiber server
func Run(addr string) error {
	// Initialize the app
	app := fiber.New()

	// Create a client manager
	manager := NewClientManager()
	manager.Start()

	// WebSocket middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Get username from query parameter
		username := c.Query("username")
		if username == "" {
			username = fmt.Sprintf("User-%d", time.Now().Unix()%10000)
		}

		// Create and register new client
		client := &Client{
			Conn:     c,
			Username: username,
		}
		manager.register <- client

		defer func() {
			manager.unregister <- c
		}()

		// Handle incoming messages
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			if messageType == websocket.TextMessage {
				// Format message with username
				formattedMsg := fmt.Sprintf("%s: %s", username, string(msg))
				// Broadcast the message to all clients
				manager.broadcast <- []byte(formattedMsg)
			}
		}
	}))

	// Serve static files
	app.Static("/", "./public")

	// Root endpoint to serve the HTML client
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	// Serve styles.css
	app.Get("/styles.css", func(c *fiber.Ctx) error {
		return c.SendFile("./public/styles.css")
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})

	// Start the server
	return app.Listen(addr)
}
