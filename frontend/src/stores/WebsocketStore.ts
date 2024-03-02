import type { VirtualCanvas } from '@/types/VirtualCanvas';
import { defineStore } from 'pinia'

export const useWebsocketStore = () => { 
  const innerStore = defineStore('websocket', {
    state: () => ({ 
      ws: undefined as WebSocket | undefined, 
      connectionError: false,
      joinedRoomId: "",
      virtualCanvas: undefined as VirtualCanvas | undefined,
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

          if (!("cmd" in jsonData)) return;

          switch (jsonData.cmd) {
            case "created_room": {
              this.joinedRoomId = jsonData.room_id;
              break;
            }
            case "joined_room": {
              this.joinedRoomId = jsonData.room_id;
              break;
            }
            case "draw_new": {
              this.virtualCanvas?.startNewLine(jsonData.id, jsonData.data)
              break;
            }
            case "draw_xy": {
              this.virtualCanvas?.drawPoint(jsonData.id, jsonData.point);
              break;
            }
            case "draw_undo": {
              this.virtualCanvas?.undoLine(jsonData.id);
              break;
            }
            case "draw_redo": {
              this.virtualCanvas?.redoLine(jsonData.id);
              break;
            }
          }      
        });
      },
      setVirtualCanvas(canvas: VirtualCanvas) {
        this.virtualCanvas = canvas;
      },
      joinRoom(roomId: String) {
        this.ws?.send(JSON.stringify({cmd: "join", room_id: roomId}))
      },
      createRoom() {
        this.ws?.send(JSON.stringify({cmd: "create"}))
      },
      sendDrawNewLine(data: { color: string, width: number }) {
        this.ws?.send(JSON.stringify({cmd: "draw_new", data}))
      },
      sendDrawPoint(point: { x: number, y: number }) {
        this.ws?.send(JSON.stringify({cmd: "draw_xy", point}))
      },
      sendDrawUndo(){
        this.ws?.send(JSON.stringify({cmd: "draw_undo"}))
      },
      sendDrawRedo(){
        this.ws?.send(JSON.stringify({cmd: "draw_redo"}))
      },
    },
  });

  const s = innerStore();
  if (s.ws == undefined) {
    s.initConnection()
  }
  return s;
}