package wordle

import "time"

type GameOptions struct {
	TurnLength time.Duration `json:"turnLength"`
	WordLength int           `json:"wordLength"`
	MaxGuesses int           `json:"maxGuesses"`
	Word       string        `json:"word,omitempty"`
}

var (
	DefaultOptions = GameOptions{
		TurnLength: 30 * time.Second,
		WordLength: 0,
		MaxGuesses: -1,
		Word:       "",
	}
)
