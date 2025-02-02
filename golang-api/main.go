package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

var clients = make(map[*websocket.Conn]bool)

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("📩 Received: %s", message)

		for client := range clients {
			if err := client.WriteMessage(messageType, message); err != nil {
				log.Println("Broadcast error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ws", handleWebSocket)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server error:", err)
	}
}