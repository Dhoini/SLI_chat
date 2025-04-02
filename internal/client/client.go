package client

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

// Run запускает консольного клиента
func Run(url string) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	// Горутина для чтения входящих сообщений
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Println("Received:", string(msg))
		}
	}()

	// Чтение ввода пользователя и отправка сообщений
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
