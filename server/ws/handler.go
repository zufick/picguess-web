package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn)

	go client.readPump()
	go client.writePump()
}

func (c *Client) readPump() {
	defer func() {
		if c.room != nil {
			c.room.unregister <- c
		}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c.messageHandler(c.conn, message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			/*
				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				for i := 0; i < n; i++ {
					w.Write([]byte{'\n'})
					w.Write(<-c.send)
				}
			*/

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) messageHandler(conn *websocket.Conn, rawMessage []byte) {
	message := map[string]string{}
	json.Unmarshal(rawMessage, &message)

	switch message["cmd"] {
	case MessageCmdCreate:
		createRoomHandler(c)
	case MessageCmdJoin:
		var data JsonMessageCmdJoin
		json.Unmarshal(rawMessage, &data)
		joinRoomHandler(c, &data)
	case MessageCmdDraw_new,
		MessageCmdDraw_xy,
		MessageCmdDraw_undo,
		MessageCmdDraw_redo:
		c.room.broadcast <- BroadcastInfo{Sender: c.id, Data: rawMessage}
	}
}

func createRoomHandler(c *Client) {
	room := newRoom()

	if err := c.JoinRoomById(room.id); err != nil {
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{ \"cmd\": \"created_room\", \"room_id\": \"%s\" }", room.id))); err != nil {
		log.Println(err)
		return
	}

}

func joinRoomHandler(c *Client, m *JsonMessageCmdJoin) {

	if err := c.JoinRoomById(m.RoomId); err != nil {
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{ \"cmd\": \"joined_room\", \"room_id\": \"%s\" }", m.RoomId))); err != nil {
		log.Println(err)
		return
	}
}
