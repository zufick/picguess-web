package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	MessageCmdCreate    = "create"
	MessageCmdJoin      = "join"
	MessageCmdDraw_new  = "draw_new"
	MessageCmdDraw_xy   = "draw_xy"
	MessageCmdDraw_redo = "draw_redo"
	MessageCmdDraw_undo = "draw_undo"
)

type JsonMessage struct {
	Cmd string `json:"cmd"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	//client := &Client{room: hub, conn: conn, send: make(chan []byte, 256)}
	//client.hub.register <- client

	for {
		var m JsonMessage

		if err := conn.ReadJSON(&m); err != nil {
			log.Println(err)
			return
		}

		go messageHandler(conn, &m)

		/*
			type message struct {
				Type string `json:"type"`
			}

			var m message
			if err := conn.ReadJSON(&m); err != nil {
				log.Println(err)
				return
			}

			if err := conn.WriteMessage(messageType, []byte("hello")); err != nil {
				log.Println(err)
				return
			}
		*/
	}

}

func messageHandler(conn *websocket.Conn, m *JsonMessage) {
	if m.Cmd == MessageCmdCreate {
		createRoomHandler(conn, m)
		return
	}
}

func createRoomHandler(conn *websocket.Conn, m *JsonMessage) {
	newRoom()
	if err := conn.WriteMessage(websocket.TextMessage, []byte("asdaadads")); err != nil {
		log.Println(err)
		return
	}

}
