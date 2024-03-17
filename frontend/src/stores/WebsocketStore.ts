import type { VirtualCanvas } from '@/types/VirtualCanvas';
import type { UserInfo } from '@/types/UserInfo';
import { defineStore } from 'pinia'
import { useGameStore } from './GameStore';
import { useRouter } from 'vue-router';
import type { VirtualCanvasDrawPoint } from '@/types/VirtualCanvas'
import type { RoomState } from '@/types/RoomState';


export const useWebsocketStore = () => { 
  const gameStore = useGameStore();
  const router = useRouter();
  let virtualCanvas: VirtualCanvas;

  const innerStore = defineStore('websocket', {
    state: () => ({ 
      ws: undefined as WebSocket | undefined, 
      connectionError: false,
      joinedRoomId: "",
      localUserId: "",
      roomState: {} as RoomState,
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

          if ("xy" in jsonData) {
            virtualCanvas?.drawPointArray(jsonData.id, jsonData.xy);
            return;
          }

          if (!("cmd" in jsonData)) return;

          switch (jsonData.cmd) {
            case "created_room": {
              this.setJoinedRoomId(jsonData.room_id, jsonData.local_user_id)
              break;
            }
            case "joined_room": {
              this.setJoinedRoomId(jsonData.room_id, jsonData.local_user_id)
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
              virtualCanvas?.startNewLine(jsonData.id, jsonData.data)
              break;
            }
            case "draw_xy": {
              virtualCanvas?.drawPoint(jsonData.id, jsonData.point);
              break;
            }
            case "draw_undo": {
              virtualCanvas?.undoLine(jsonData.id);
              break;
            }
            case "draw_redo": {
              virtualCanvas?.redoLine(jsonData.id);
              break;
            }
            case "draw_clear": {
              virtualCanvas?.drawClear(jsonData.id);
              break;
            }
          }
        });
      },
      setVirtualCanvas(canvas: VirtualCanvas) {
        virtualCanvas = canvas;
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
      sendDrawPoints(points: number[]) {
        this.ws?.send(JSON.stringify({xy: points}))
      },
      sendDrawUndo(){
        this.ws?.send(JSON.stringify({cmd: "draw_undo"}))
      },
      sendDrawRedo(){
        this.ws?.send(JSON.stringify({cmd: "draw_redo"}))
      },
      sendDrawClear(voted: boolean) {
        this.ws?.send(JSON.stringify({cmd: "draw_clear", voted}))
      },
      startGame() {
        this.ws?.send(JSON.stringify({cmd: "game_start"}))
      },
      setJoinedRoomId(id: string, localUserId: string) {
        this.joinedRoomId = id;
        this.localUserId = localUserId
        router.push(`/room/${id}`)
      },
      pickAnswer(answer: string) {
        this.ws?.send(JSON.stringify({cmd: "game_answer", answer}))
      }
    },
  });

  const s = innerStore();
  if (s.ws == undefined) {
    s.initConnection()
  }
  return s;
}