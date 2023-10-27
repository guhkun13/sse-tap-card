// main.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const CONST_CARD_UUID = "CARD-UUID"
const CONST_NPM = "NPM"

const MessageTypePing = "ping"
const MessageTypeData = "data"

var globalChan chan string

type Message struct {
	Type     string `json:"type"`
	CardUUID string `json:"cardUUID"`
	Text     string `json:"text"`
	// NPM      string `json:"npm"`

}

type CardData struct {
	Name    string
	CardId  string
	Faculty string
	OrgCode string
}

type CardResponse struct {
	IsValid     bool        `json:"isValid"`
	HasRegister bool        `json:"hasRegister"`
	Data        interface{} `json:"data"`
}

var cardList = []CardData{
	{Name: "Abdul Amin", CardId: "123", Faculty: "Engineering", OrgCode: "E001"},
	{Name: "Budi Anduk", CardId: "456", Faculty: "Science", OrgCode: "S002"},
	{Name: "Dodit", CardId: "789", Faculty: "Law", OrgCode: "L012"},
	{Name: "Radit", CardId: "111", Faculty: "Computer Science", OrgCode: "CS002"},
	{Name: "Kevin", CardId: "222", Faculty: "Math", OrgCode: "M001"},
}

var registeredCardList = []CardData{
	{Name: "Abdul Amin", CardId: "123", Faculty: "Engineering", OrgCode: "E001"},
	{Name: "Radit", CardId: "111", Faculty: "Computer Science", OrgCode: "CS002"},
}

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
			fmt.Println("send data to client")
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
	cardUUID := chi.URLParam(r, "id")
	fmt.Println(cardUUID)
	if cardUUID != "" {
		fmt.Println("assign to globalChan")
		globalChan <- cardUUID
	}
}

func getCardHandler(w http.ResponseWriter, r *http.Request) {

	cardId := chi.URLParam(r, "id")
	fmt.Println("getCardHandler", cardId)

	// validate npm exist
	cardData := getCardById(cardId)
	if cardData == nil {
		fmt.Println("case 1: Card not exist")
		returnResponse(w, CardResponse{
			IsValid:     false,
			HasRegister: false,
			Data:        cardData,
		})
		return
	}

	registeredCard := getRegisteredCardById(cardId)
	if registeredCard != nil {
		fmt.Println("case 2: Card has registered")
		returnResponse(w, CardResponse{
			IsValid:     true,
			HasRegister: true,
			Data:        registeredCard,
		})
		return
	}

	fmt.Println("case 3: Card valid")
	returnResponse(w, CardResponse{
		IsValid:     true,
		HasRegister: false,
		Data:        cardData,
	})
}

func getCardById(cardId string) (cardData interface{}) {
	for _, card := range cardList {
		if card.CardId == cardId {
			cardData = card
			break
		}
	}

	return cardData
}

func getRegisteredCardById(cardId string) (cardData interface{}) {
	for _, card := range registeredCardList {
		if card.CardId == cardId {
			cardData = card
			break
		}
	}

	return cardData
}

func returnResponse(w http.ResponseWriter, cardReponse CardResponse) {
	fmt.Println("return response")

	response := map[string]interface{}{
		"isValid":     cardReponse.IsValid,
		"hasRegister": cardReponse.HasRegister,
		"data":        cardReponse.Data,
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// make global channel
	globalChan = make(chan string)

	r := chi.NewRouter()
	r.Get("/sse", handleSSE)
	r.Get("/tap-card/{id}", handleTapCard)
	r.Get("/get-card/{id}", getCardHandler)

	fmt.Println("Server started on :8082")
	err := http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatal(err)
	}
}
