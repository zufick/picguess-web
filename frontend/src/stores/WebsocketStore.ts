import { defineStore } from 'pinia'

export const useWebsocketStore = () => { 
  const innerStore = defineStore('websocket', {
    state: () => ({ ws: undefined, connectionError: false }),
    actions: {
      initConnection() {
        console.log("Initializing ws connection...")

        let socket = new WebSocket("ws://localhost:8080");
        this.ws = socket;

        socket.addEventListener("error", (event) => {
            this.connectionError = true
        });
        
        socket.addEventListener("close", (event) => {
          this.connectionError = true
        });
        
        // Listen for messages
        socket.addEventListener("message", (event) => {
            console.log("Message from server ", event.data);
        });

      },
    },
  });

  const s = innerStore();
  if (s.ws == undefined) {
    s.initConnection()
  }
  return s;
}