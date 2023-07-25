package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/messages", messageHandler).Methods("POST")
	router.HandleFunc("/ws", handleConnections)

	go handleMessages()

	http.ListenAndServe(":8080", router)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Handle registration
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Handle login and JWT creation
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	// Handle sending of messages
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
