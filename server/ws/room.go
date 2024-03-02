package ws

import (
	"log"

	"github.com/google/uuid"
)

type Room struct {
	id string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var rooms = make(map[string]*Room)

func newRoom() *Room {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
	}

	room := &Room{
		id:         id.String(),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	room.run()
	rooms[room.id] = room
	return room
}

func (h *Room) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
