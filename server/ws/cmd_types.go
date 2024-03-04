package ws

const (
	MessageCmdCreate    = "create"
	MessageCmdJoin      = "join"
	MessageCmdDraw_new  = "draw_new"
	MessageCmdDraw_xy   = "draw_xy"
	MessageCmdDraw_redo = "draw_redo"
	MessageCmdDraw_undo = "draw_undo"
)

type JsonMessageCmd struct {
	Cmd string `json:"cmd"`
}

type JsonMessageCmdJoin struct {
	Cmd    string `json:"cmd"`
	RoomId string `json:"room_id"`
}

// Draws a new line
type JsonMessageCmdDraw_new struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Color string `json:"color"`
		Width string `json:"width"`
	} `json:"data"`
}

type JsonMessageCmdDraw_xy struct {
	Cmd   string `json:"cmd"`
	Point struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"point"`
}
