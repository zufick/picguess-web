import { WebSocketServer } from "ws";
import { v4 as uuidv4 } from "uuid";

const wss = new WebSocketServer({ port: 8080 });

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
      case "join":
        if (!rooms[message.room]) {
          ws.send(JSON.stringify({ error: "room not found" }));
          break;
        }

        rooms[message.room][connUuid] = ws;
        ws.send(JSON.stringify({ message: "successfully joined room" }));
        break;
      case "create":
        let roomId = "room" + uuidv4();
        rooms[roomId] = { [connUuid]: ws };
        console.log(rooms);
        ws.send(JSON.stringify({ roomId: roomId }));
        break;
    }
  });
});
