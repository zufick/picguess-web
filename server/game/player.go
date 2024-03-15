package game

type Player struct {
	Score            int             `json:"score"`
	WordPool         []string        `json:"wordPool"`
	OpponentWinWords []string        `json:"opponentWinWords"`
	AnswerResults    []AnswerResults `json:"answerResults"`
	OpponentResults  []AnswerResults `json:"opponentResults"`
	Opponent         *Player         `json:"-"`
}

func (p *Player) GetCorrectAnswers() []AnswerResults {
	var result []AnswerResults
	for _, a := range p.AnswerResults {
		if a.IsCorrect {
			result = append(result, a)
		}
	}
	return result
}
