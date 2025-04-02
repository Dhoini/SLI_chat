package server

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Client представляет подключенного WebSocket-клиента
type Client struct {
	Conn     *websocket.Conn
	Username string
}

// ClientManager управляет WebSocket-клиентами и рассылкой
type ClientManager struct {
	clients    map[*websocket.Conn]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

// NewClientManager создает нового менеджера клиентов
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*websocket.Conn]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *websocket.Conn),
	}
}

// Start начинает обработку управления клиентами
func (m *ClientManager) Start() {
	go func() {
		for {
			select {
			case client := <-m.register:
				m.mutex.Lock()
				m.clients[client.Conn] = client
				m.mutex.Unlock()
				log.Printf("Клиент подключен: %s", client.Username)

				// Оповещение всех клиентов о новом пользователе
				joinMsg := fmt.Sprintf("Система: %s присоединился к чату", client.Username)
				m.broadcast <- []byte(joinMsg)

			case conn := <-m.unregister:
				m.mutex.Lock()
				if client, ok := m.clients[conn]; ok {
					username := client.Username
					delete(m.clients, conn)
					conn.Close()

					// Оповещение о выходе пользователя
					leaveMsg := fmt.Sprintf("Система: %s покинул чат", username)
					m.broadcast <- []byte(leaveMsg)
				}
				m.mutex.Unlock()
				log.Println("Клиент отключен")

			case message := <-m.broadcast:
				m.mutex.Lock()
				clientCount := len(m.clients)
				log.Printf("Рассылка %d клиентам: %s", clientCount, string(message))
				for conn, _ := range m.clients {
					if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
						log.Println("Ошибка отправки:", err)
						conn.Close()
						delete(m.clients, conn)
					}
				}
				m.mutex.Unlock()
			}
		}
	}()
}

// Run запускает Fiber-сервер
func Run(addr string) error {
	// Инициализация приложения
	app := fiber.New()

	// Создание менеджера клиентов
	manager := NewClientManager()
	manager.Start()

	// Middleware для WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket-эндпоинт
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Получение имени пользователя из параметров
		username := c.Query("username")
		if username == "" {
			username = fmt.Sprintf("Пользователь-%d", time.Now().Unix()%10000)
		}

		// Создание и регистрация нового клиента
		client := &Client{
			Conn:     c,
			Username: username,
		}
		manager.register <- client

		defer func() {
			manager.unregister <- c
		}()

		// Обработка входящих сообщений
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения:", err)
				break
			}

			if messageType == websocket.TextMessage {
				// Форматирование сообщения с именем пользователя
				formattedMsg := fmt.Sprintf("%s: %s", username, string(msg))
				log.Printf("Рассылка сообщения: %s", formattedMsg)
				// Рассылка сообщения всем клиентам
				manager.broadcast <- []byte(formattedMsg)
			}
		}
	}))

	// Старт сервера
	return app.Listen(addr)
}
