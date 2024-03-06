package game

import (
	"game-server/utils/lang"
	"log"
)

type Player struct {
	Score            int          `json:"score"`
	WordPool         []string     `json:"wordPool"`
	OpponentWinWords []string     `json:"opponentWinWords"`
	AnswerResults    map[int]bool `json:"answerResults"`
	Opponent         *Player      `json:"-"`
}

type Game struct {
	players        map[*Player]bool
	wordsPerPlayer int
	wordsToWin     int
	StateUpdated   chan int
}

func NewGame() *Game {
	return &Game{
		players:        make(map[*Player]bool),
		wordsPerPlayer: 16,
		wordsToWin:     5,
	}
}

func (g *Game) NewPlayer() *Player {
	p := &Player{
		Score: 0,
	}
	g.players[p] = true
	return p
}

func (g *Game) Start() {
	g.StateUpdated = make(chan int)
	words, err := lang.GenerateWords("ru", g.wordsPerPlayer*len(g.players))

	if err != nil {
		log.Println(err)
		return
	}

	var players []*Player
	k := 0
	for p := range g.players {
		players = append(players, p)
		p.WordPool = words[g.wordsPerPlayer*(k) : g.wordsPerPlayer*(k+1)]
		k++
	}

	for i := 0; i < len(players); i++ {
		players[i].Opponent = players[(len(players)-1)-i]
		players[i].OpponentWinWords = lang.GetRandomWords(players[i].Opponent.WordPool, g.wordsToWin)

	}

	g.StateUpdated <- 1
}
