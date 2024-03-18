package game

import (
	"game-server/utils/lang"
	"log"

	"golang.org/x/exp/slices"
)

type AnswerResults struct {
	Answer    string `json:"answer"`
	IsCorrect bool   `json:"isCorrect"`
}

type Game struct {
	players        map[*Player]bool
	wordsPerPlayer int
	wordsToWin     int
	winner         *Player
	StateUpdated   chan int
	lang           string
}

func NewGame() *Game {
	return &Game{
		players:        make(map[*Player]bool),
		wordsPerPlayer: 16,
		wordsToWin:     5,
		lang:           "ru",
	}
}

func (g *Game) NewPlayer(id int) *Player {
	p := &Player{
		Id:            id,
		Score:         0,
		AnswerResults: make([]AnswerResults, 0),
	}
	g.players[p] = true
	return p
}

func (g *Game) Start() {
	g.StateUpdated = make(chan int)
	words, err := lang.GenerateWords(g.lang, g.wordsPerPlayer*len(g.players))

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

func (g *Game) AnswerWord(p *Player, answer string) {
	if p == nil || len(p.AnswerResults) >= g.wordsPerPlayer || len(p.GetCorrectAnswers()) >= g.wordsToWin {
		return
	}

	for _, a := range p.AnswerResults {
		if a.Answer == answer {
			return
		}
	}

	isCorrect := slices.Contains(p.Opponent.OpponentWinWords, answer)

	p.AnswerResults = append(p.AnswerResults, AnswerResults{Answer: answer, IsCorrect: isCorrect})
	p.Opponent.OpponentResults = p.AnswerResults

	if isCorrect {
		p.Score += 1
		p.Opponent.Score += 2
	} else {
		p.Score -= 2
		p.Opponent.Score -= 1
	}

	if p.Score < 0 {
		p.Score = 0
	}

	if p.Opponent.Score < 0 {
		p.Opponent.Score = 0
	}

	g.CheckWin()

	g.StateUpdated <- 1
}

func (g *Game) CheckWin() {
	wordsRequired := len(g.players) * g.wordsToWin
	correctAnsweredWords := 0

	var playerWithMostScore *Player

	for p := range g.players {

		for _, a := range p.AnswerResults {
			if a.IsCorrect {
				correctAnsweredWords++
			}
		}

		if playerWithMostScore == nil {
			playerWithMostScore = p
		}

		if playerWithMostScore.Score < p.Score {
			playerWithMostScore = p
		}
	}

	if correctAnsweredWords >= wordsRequired {
		g.winner = playerWithMostScore
	}
}

func (g *Game) GetWinner() *Player {
	return g.winner
}

func (g *Game) GetPlayers() []*Player {
	var players []*Player

	for p := range g.players {
		players = append(players, p)
	}

	return players
}
