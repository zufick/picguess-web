package ws

import (
	"errors"
	"game-server/game"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     int
	room   *Room
	player *game.Player

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	userInfo UserInfo
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		id:   0,
		room: nil,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) JoinRoomById(roomId string, userInfo UserInfo) error {
	if rooms[roomId] == nil {
		return errors.New("Room not found")
	}
	lastUserIdInRoom++
	c.id = lastUserIdInRoom
	c.userInfo = userInfo

	c.room = rooms[roomId]
	c.room.register <- c

	return nil
}
