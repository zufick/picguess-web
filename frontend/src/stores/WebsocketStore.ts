import { defineStore } from 'pinia'

export const useWebsocketStore = () => { 
  const innerStore = defineStore('websocket', {
    state: () => ({ 
      ws: undefined as WebSocket | undefined, 
      connectionError: false,
      joinedRoomId: "",
      canvasRedraw: undefined as Function | undefined,
      canvasDrawStop: undefined as Function | undefined,
    }),
    getters: {
      isJoinedRoom(): Boolean {
        return !this.connectionError && this.joinedRoomId != "";
      }
    },
    actions: {
      initConnection() {
        console.log("Initializing ws connection...")

        let socket = new WebSocket("ws://" + location.hostname + ":8081");
        this.ws = socket;

        socket.addEventListener("error", (event) => {
            this.connectionError = true
        });
        
        socket.addEventListener("close", (event) => {
          this.connectionError = true
        });
        
        // Listen for messages
        socket.addEventListener("message", (event) => {
          let jsonData;
          try{
            jsonData = JSON.parse(event.data);
          } catch(e) {
            console.error(e);
            return;
          }

          if (!("type" in jsonData)) return;

          switch (jsonData.type) {
            case "created_room": {
              this.joinedRoomId = jsonData.room_id;
              break;
            }
            case "joined_room": {
              this.joinedRoomId = jsonData.room_id;
              break;
            }
            case "draw_new": {
              
              break;
            }
            case "draw": {
              this.canvasRedraw?.(jsonData.data);
              break;
            }
            case "drawstop": {
              this.canvasDrawStop?.();
              break;
            }
          }      
        });
      },
      setCanvasFunctions(redraw: any, stop: any) {
        this.canvasRedraw = redraw;
        this.canvasDrawStop = stop;
      },
      joinRoom(roomId: String) {
        this.ws?.send(JSON.stringify({type: "join", room_id: roomId}))
      },
      createRoom() {
        this.ws?.send(JSON.stringify({type: "create"}))
      },
      sendDrawNewLine(data: any) {
        this.ws?.send(JSON.stringify({type: "draw_new", data}))
      },
      sendDraw(data: any) {
        this.ws?.send(JSON.stringify({type: "draw", data}))
      },
      sendDrawStop() {
        this.ws?.send(JSON.stringify({type: "drawstop"}))
      }
    },
  });

  const s = innerStore();
  if (s.ws == undefined) {
    s.initConnection()
  }
  return s;
}