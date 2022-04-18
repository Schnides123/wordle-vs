package wordle

import "time"

type Guess struct {
	Word      string    `json:"word"`
	Result    []int     `json:"result"`
	Player    string    `json:"player"`
	Timestamp time.Time `json:"timestamp"`
}

func NewGuess(word string, result []int, player string) *Guess {
	return &Guess{
		Word:      word,
		Result:    result,
		Player:    player,
		Timestamp: time.Now(),
	}
}
