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
	// Запрос имени пользователя
	var username string
	fmt.Print("Введите ваше имя: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
	}

	// Добавление имени пользователя к URL WebSocket
	if !strings.Contains(url, "?") {
		url = url + "?username=" + username
	} else {
		url = url + "&username=" + username
	}

	fmt.Printf("Подключение к %s...\n", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу")
	fmt.Println("Введите сообщение и нажмите Enter для отправки. Введите 'exit' для выхода.")

	// Горутина для чтения входящих сообщений
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Соединение закрыто:", err)
				os.Exit(0)
				return
			}
			fmt.Printf("\n%s\n> ", string(msg))
		}
	}()

	// Чтение ввода пользователя и отправка сообщений
	scanner = bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Выход...")
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			log.Println("Ошибка отправки:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Ошибка ввода:", err)
	}
}
