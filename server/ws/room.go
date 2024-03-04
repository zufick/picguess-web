package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type BroadcastInfo struct {
	Sender int
	Data   []byte
}

type Room struct {
	id string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan BroadcastInfo

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var rooms = make(map[string]*Room)
var lastUserIdInRoom = 0

func newRoom() *Room {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
	}

	room := &Room{
		id:         id.String(),
		broadcast:  make(chan BroadcastInfo),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go room.run()
	rooms[room.id] = room
	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast:
			var messageMap map[string]json.RawMessage
			err := json.Unmarshal(message.Data, &messageMap)
			if err != nil {
				fmt.Println(err)
			}
			messageMap["id"] = json.RawMessage(fmt.Sprintf("\"%d\"", message.Sender))

			messageWithSender, err := json.Marshal(&messageMap)
			if err != nil {
				fmt.Println(err)
			}

			for client := range r.clients {

				if client.id == message.Sender {
					continue
				}

				select {
				case client.send <- messageWithSender:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
