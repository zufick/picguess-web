import type { VirtualCanvas } from '@/types/VirtualCanvas';
import type { UserInfo } from '@/types/UserInfo';
import { defineStore } from 'pinia'
import { useGameStore } from './GameStore';
import { useRouter } from 'vue-router';


export const useWebsocketStore = () => { 
  const gameStore = useGameStore();
  const router = useRouter();

  const innerStore = defineStore('websocket', {
    state: () => ({ 
      ws: undefined as WebSocket | undefined, 
      connectionError: false,
      joinedRoomId: "",
      roomState: {},
      virtualCanvas: undefined as VirtualCanvas | undefined,
    }),
    getters: {
      isConnected() : Boolean {
        return this.ws != undefined
      },
      isJoinedRoom(): Boolean {
        return !this.connectionError && this.joinedRoomId != "";
      }
    },
    actions: {
      initConnection() {
        console.log("Initializing ws connection...")

        let socket = new WebSocket("ws://" + location.hostname + ":8081");

        socket.addEventListener("open", (event) => {
          this.ws = socket;
        });

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
              this.setJoinedRoomId(jsonData.room_id)
              break;
            }
            case "joined_room": {
              this.setJoinedRoomId(jsonData.room_id)
              break;
            }
            case "roomstate": {
              this.roomState = jsonData;
              break;
            }
            case "gamestate": {
              gameStore.gameState = jsonData.gamestate
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
      joinRoom(roomId: String, userInfo: UserInfo) {
        this.ws?.send(JSON.stringify({cmd: "join", room_id: roomId, user_info: userInfo}))
      },
      createRoom(userInfo: UserInfo) {
        this.ws?.send(JSON.stringify({cmd: "create", user_info: userInfo}))
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
      startGame() {
        this.ws?.send(JSON.stringify({cmd: "game_start"}))
      },
      setJoinedRoomId(id: string) {
        this.joinedRoomId = id;
        router.push(`/room/${id}`)
      }
    },
  });

  const s = innerStore();
  if (s.ws == undefined) {
    s.initConnection()
  }
  return s;
}