package wordle

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/util"

	"github.com/google/uuid"
)

type Game struct {
	Id      string      `json:"id"`
	Word    *Word       `json:"word"`
	Guessed []*Guess    `json:"guessed"`
	Winner  bool        `json:"winner"`
	Options GameOptions `json:"options"`
	Players *util.Set   `json:"players"`

	playersLock sync.Mutex
}

func NewGame(options GameOptions) *Game {

	word := options.Word
	if len(options.Word) > 0 {
		options.WordLength = len(options.Word)
	} else {
		word = randomWord(options.WordLength)
	}
	options.Word = ""
	return &Game{
		Id:          uuid.NewString(),
		Word:        NewWord(word),
		Guessed:     []*Guess{},
		Winner:      false,
		Options:     options,
		Players:     util.NewSet(),
		playersLock: sync.Mutex{},
	}
}

func (g *Game) Guess(word, player string) error {
	word = strings.ToLower(word)
	if g.IsCompleted() {
		return errors.New("game is completed")
	} else if len(word) != g.Word.Length {
		return errors.New("invalid word length")
	}
	if len(g.Guessed) > 0 {
		lastGuess := g.Guessed[len(g.Guessed)-1]
		if lastGuess.Player == player && time.Now().Before(lastGuess.Timestamp.Add(g.Options.TurnLength)) {
			return errors.New("You have to wait at least 30 seconds before stealing the opponent's turn")
		}
	}

	result := g.Word.Check(word)
	flag := true
	for _, v := range result {
		if v != 2 {
			flag = false
			break
		}
	}
	g.Winner = flag
	guess := NewGuess(word, result, player)
	g.Guessed = append(g.Guessed, guess)
	return nil
}

func (g *Game) IsCompleted() bool {
	return g.Winner || len(g.Guessed) == g.Options.MaxGuesses
}

func (g *Game) AddPlayer(id string) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()

	g.Players.Add(id)
}

func (g *Game) RemovePlayer(id string) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()

	g.Players.Remove(id)
}

func (g *Game) Cancer() {
	g.playersLock.Lock()
}
func (g *Game) UnCancer() {
	g.playersLock.Unlock()
}
