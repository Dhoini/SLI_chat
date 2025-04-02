package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// Run запускает консольного клиента
func Run(url string) {
	// Ask for username
	var username string
	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
	}

	// Add username to websocket URL
	if !strings.Contains(url, "?") {
		url = url + "?username=" + username
	} else {
		url = url + "&username=" + username
	}

	fmt.Printf("Connecting to %s...\n", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")
	fmt.Println("Type your message and press Enter to send. Type 'exit' to quit.")

	// Горутина для чтения входящих сообщений
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Connection closed:", err)
				os.Exit(0)
				return
			}
			fmt.Printf("Received: %s\n", string(msg))
		}
	}()

	// Чтение ввода пользователя и отправка сообщений
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Exiting...")
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
