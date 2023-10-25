// main.go

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var origins = []string{"http://localhost", "http://localhost:8082", "http://127.0.0.1:8082", "127.0.0.1:8082", "http://localhost:8082", "http://localhost:5173"}

// var origins = []string{"*"}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// var origin = r.Header.Get("origin")
		// for _, allowOrigin := range origins {
		// 	if origin == allowOrigin {
		// 		return true
		// 	}
		// }
		// return false
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from client
		fmt.Println("readMessage!")
		messageType, message, err := conn.ReadMessage()
		fmt.Println("messageType = ", messageType)

		if err != nil {
			log.Println("WebSocket read error:", err)
			// break
		}

		// Print received message
		log.Printf("Received message: %s", message)
		log.Println("then SEND ")

		// Send a response to the client
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

func handleEventWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleEventWebSocket")
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from client
		fmt.Println("readMessage!")
		messageType, message, err := conn.ReadMessage()
		fmt.Println("messageType!", messageType)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// Print received message
		log.Printf("Received message: %s", message)
		log.Println("then SEND ")

		// resMessage := fmt.Sprintf("[ok] %s", message)
		// fmt.Println("prepared message: ", resMessage)

		// Send a response to the client
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
		log.Println("done writing message")
	}
}

func main2() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/event", handleEventWebSocket)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
