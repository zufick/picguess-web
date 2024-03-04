package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int
	room *Room

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		id:   0,
		room: nil,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) JoinRoomById(roomId string) error {
	if rooms[roomId] == nil {
		return errors.New("Room not found")
	}

	c.room = rooms[roomId]
	c.room.register <- c

	lastUserIdInRoom++
	c.id = lastUserIdInRoom

	return nil
}
