package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const PORT = ":8090"

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

		type message struct {
			Type string `json:"type"`
		}

		var m message
		if err := conn.ReadJSON(&m); err != nil {
			log.Println(err)
			return
		}

		log.Println("aa")

		if err := conn.WriteMessage(messageType, []byte("hello")); err != nil {
			log.Println(err)
			return
		}
	}

}

func main() {
	fmt.Printf("Starting http server at %s\n", PORT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	httpServer := &http.Server{
		Addr: PORT,
	}

	http.HandleFunc("/ws", handler)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
