package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcastChan = make(chan string)

var upgrader = websocket.Upgrader{}

func Run(addr string) error {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	return http.ListenAndServe(addr, nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("Client connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, ws)
			break
		}
		broadcastChan <- string(msg)
	}
}
