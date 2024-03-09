package ws

import "game-server/game"

const (
	MessageCmdCreate      = "create"
	MessageCmdJoin        = "join"
	MessageCmdGame_start  = "game_start"
	MessageCmdGame_answer = "game_answer"
	MessageCmdDraw_new    = "draw_new"
	MessageCmdDraw_xy     = "draw_xy"
	MessageCmdDraw_redo   = "draw_redo"
	MessageCmdDraw_undo   = "draw_undo"
	MessageCmdDraw_clear  = "draw_clear"
	MessageCmdGameState   = "gamestate"
)

type UserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type JsonMessageCmd struct {
	Cmd string `json:"cmd"`
}

type JsonMessageCmdCreate struct {
	Cmd      string   `json:"cmd"`
	UserInfo UserInfo `json:"user_info"`
}

type JsonMessageCmdJoin struct {
	Cmd      string   `json:"cmd"`
	RoomId   string   `json:"room_id"`
	UserInfo UserInfo `json:"user_info"`
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

type JsonMessageCmdGame_answer struct {
	Cmd    string `json:"cmd"`
	Answer string `json:"answer"`
}

type RoomStateBroadcastData struct {
	Cmd         string     `json:"cmd"`
	ClientInfos []UserInfo `json:"clients"`
}

type GameStateBroadcast struct {
	Cmd   string `json:"cmd"`
	State struct {
		Player *game.Player              `json:"player"`
		Winner *GameStateBroadcastWinner `json:"winner"`
	} `json:"gamestate"`
}

type GameStateBroadcastWinner struct {
	UserInfo UserInfo `json:"user_info"`
	Score    int      `json:"score"`
}
