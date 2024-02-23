import { WebSocketServer } from "ws";
import { v4 as uuidv4 } from "uuid";


const wss = new WebSocketServer({ port: 8081 });

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

        rooms[roomId][connUuid] = ws;
        ws.send(JSON.stringify({ type: "joined_room", "room_id": roomId }));
        break;
      }
      case "create": {
        let roomId = "room" + uuidv4();
        rooms[roomId] = { [connUuid]: ws };
        ws.send(JSON.stringify({ type: "created_room", "room_id": roomId }));
        break;
      }
    }
  });
});