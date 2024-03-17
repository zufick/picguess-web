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

	clearCanvasVotes map[*Client]bool

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
		clearCanvasVotes:    make(map[*Client]bool),
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
			go r.broadcastRoomState()
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			go r.broadcastRoomState()
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

func (r *Room) VoteClearCanvas(c *Client, voted bool) {
	r.clearCanvasVotes[c] = voted
	firstVoterId := -1

	if len(r.clearCanvasVotes) >= len(r.clients) {
		inFavor := 0
		for c, voted := range r.clearCanvasVotes {
			if firstVoterId == -1 {
				firstVoterId = c.id
			}

			if voted {
				inFavor++
			}
		}

		if len(r.clearCanvasVotes) >= len(r.clients) {
			r.clearCanvasVotes = make(map[*Client]bool)

			if inFavor > len(r.clients)/2 {
				c.room.broadcast <- []byte("{ \"cmd\": \"draw_clear\" }")
			}
		}
	}

	go r.broadcastRoomState()
}

func (r *Room) GetClearCanvasVoteResults() []bool {
	var results []bool
	for _, v := range r.clearCanvasVotes {
		results = append(results, v)
	}

	fmt.Println(results)
	return results
}

func (r *Room) broadcastRoomState() {
	var clientsInfo []UserInfo

	for client := range r.clients {
		clientsInfo = append(clientsInfo, client.userInfo)
	}

	clearCanvasVotes := make(map[int]bool)
	for c, v := range r.clearCanvasVotes {
		clearCanvasVotes[c.id] = v
	}

	roomState := &RoomStateBroadcastData{
		Cmd:             "roomstate",
		ClientInfos:     clientsInfo,
		VoteClearCanvas: clearCanvasVotes,
	}

	roomStateJson, err := json.Marshal(roomState)
	if err != nil {
		fmt.Println(err)
		return
	}

	r.broadcast <- roomStateJson
}
