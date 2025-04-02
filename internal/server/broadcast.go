package server

import (
	"github.com/gorilla/websocket"
	"log"
)

// handleMessages рассылает сообщения всем клиентам
func handleMessages() {
	for {
		msg := <-broadcastChan
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
