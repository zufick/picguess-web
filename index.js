import { WebSocketServer } from "ws";
import { v4 as uuidv4 } from "uuid";

const PORT = 8081;

console.log("Starting WebSocketServer at :" + PORT)
const wss = new WebSocketServer({ port: PORT });

const rooms = {};

wss.on("connection", function connection(ws) {
  const connUuid = "user" + uuidv4();

  ws.on("error", console.error);

  ws.on("message", function message(data) {
    let message;
    try {
      message = JSON.parse(data);
    } catch (e) {
      console.error(e);
      return;
    }
    if (!("type" in message)) return;

    switch (message.type) {
      case "join": {
        let roomId = message.room_id;

        if (!rooms[roomId]) {
          ws.send(JSON.stringify({ error: "room not found" }));
          break;
        }

        rooms[roomId]['users'].push({id: connUuid, socket: ws});
        ws.roomId = roomId;
        ws.send(JSON.stringify({ type: "joined_room", "room_id": roomId }));
        break;
      }
      case "create": {
        let roomId = "room" + uuidv4();
        rooms[roomId] = { users: [ {id: connUuid, socket: ws} ] };
        ws.roomId = roomId;
        ws.send(JSON.stringify({ type: "created_room", "room_id": roomId }));
        break;
      }
      case "draw_new": {
        if(!ws.roomId) { return; }
        let room = rooms[ws.roomId];

        room.users.forEach((element, index) => {
          if(element.id == connUuid){
            return;
          }
          message.id = index;
          element.socket.send(JSON.stringify(message), { binary: false })
        });
        break;
      }
      case "draw_xy": {
        if(!ws.roomId) { return; }
        let room = rooms[ws.roomId];

        room.users.forEach((element, index) => {
          if(element.id == connUuid){
            return;
          }
          message.id = index;
          element.socket.send(JSON.stringify(message), { binary: false })
        });
        break;
      }
    }
  });
});