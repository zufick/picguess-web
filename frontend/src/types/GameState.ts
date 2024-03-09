export type GameState = {
    player: {
        score: 0,
        opponentWinWords: [],
        wordPool: [],
        opponentResults: [
            { answer: string, isCorrect: boolean }
        ]
        answerResults: [
            { answer: string, isCorrect: boolean }
        ],
    },
    winner: {
        
    }
}