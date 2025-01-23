package application

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var connections = struct {
	sync.Mutex
	clients map[*websocket.Conn]bool
}{clients: make(map[*websocket.Conn]bool)}

func StartWebSocketServer() {
	http.HandleFunc("/ws", handleWebSocket)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	connections.Lock()
	connections.clients[conn] = true
	connections.Unlock()

	log.Println("WebSocket connection established")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		broadcastMessage(messageType, message)
	}

	connections.Lock()
	delete(connections.clients, conn)
	connections.Unlock()
	log.Println("WebSocket connection closed")
}

func broadcastMessage(messageType int, message []byte) {
	connections.Lock()
	defer connections.Unlock()
	for conn := range connections.clients {
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
		}
	}
}