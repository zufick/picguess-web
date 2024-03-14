package ws

import (
	"encoding/json"
	"fmt"
	"game-server/game"
	"log"

	"github.com/google/uuid"
)

type BroadcastSenderInfo struct {
	Sender int
	Data   []byte
}

type Room struct {
	id string

	game *game.Game

	lastUserIdInRoom int

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcastFromClient chan BroadcastSenderInfo

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
		id:                  id.String(),
		lastUserIdInRoom:    0,
		broadcastFromClient: make(chan BroadcastSenderInfo),
		broadcast:           make(chan []byte),
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		clients:             make(map[*Client]bool),
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
			go r.broadcastGameState()
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			go r.broadcastGameState()
		case message := <-r.broadcastFromClient:
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
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

func (r *Room) broadcastGameState() {
	var clientsInfo []UserInfo

	for client := range r.clients {
		clientsInfo = append(clientsInfo, client.userInfo)
	}

	roomState := &RoomStateBroadcastData{
		Cmd:         "roomstate",
		ClientInfos: clientsInfo,
	}

	roomStateJson, err := json.Marshal(roomState)
	if err != nil {
		fmt.Println(err)
		return
	}

	r.broadcast <- roomStateJson
}
