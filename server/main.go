package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		type kek struct {
			name string
			age  int
		}

		var test kek
		if err := conn.ReadJSON(&test); err != nil {
			log.Println(err)
			return
		}

		log.Println(test)

		if err := conn.WriteMessage(messageType, []byte("hello")); err != nil {
			log.Println(err)
			return
		}
	}

}

func main() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8090", nil)
}
