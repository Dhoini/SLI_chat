package main

import (
	"SLI_chat/internal/server"
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

	// Читаем адрес сервера из переменной окружения
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080" // Значение по умолчанию
	}

	log.Printf("Starting Fiber server on %s...", addr)
	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}
