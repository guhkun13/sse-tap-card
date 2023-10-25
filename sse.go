// main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const CONST_CARD_UUID = "CARD-UUID"
const CONST_NPM = "NPM"

const MessageTypePing = "ping"
const MessageTypeData = "data"

type Message struct {
	Type     string `json:"type"`
	CardUUID string `json:"cardUUID"`
	Text     string `json:"text"`
	// NPM      string `json:"npm"`

}

var globalChan chan string

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send events to the client periodically
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// keep connection alive
	for {
		select {
		case cardUUID := <-globalChan:
			msg := Message{
				Type:     MessageTypeData,
				Text:     CONST_CARD_UUID,
				CardUUID: cardUUID,
			}

			jsonData, _ := json.Marshal(msg)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			w.(http.Flusher).Flush()

		case <-ticker.C:
			message := Message{
				Type: MessageTypePing,
				Text: "ping",
			}
			jsonData, _ := json.Marshal(message)

			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func handleTapCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleTapCard!")
	cardUUID := r.URL.Query().Get("cardUUID")
	if cardUUID != "" {
		globalChan <- cardUUID
	}
}

func main() {

	// make global channel
	globalChan = make(chan string)

	// define handler
	http.HandleFunc("/sse", handleSSE)
	http.HandleFunc("/tap-card", handleTapCard)

	fmt.Println("Server started on :8082")
	http.ListenAndServe(":8082", nil)
}
