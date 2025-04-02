package main

import (
	"SLI_chat/internal/client"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file, using default values")
	}

	// Читаем WebSocket URL из переменной окружения
	url := os.Getenv("WS_URL")
	if url == "" {
		url = "ws://localhost:8080/ws" // Значение по умолчанию
	}

	client.Run(url)
}
