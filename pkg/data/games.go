package data

import (
	"github.com/Schnides123/wordle-vs/pkg/wordle"
)

var (
	Games = map[string]*wordle.Game{}
)

func GetGame(gameID string) (*wordle.Game, error) {
	game, ok := Games[gameID]
	if ok {
		return game, nil
	}
	return nil, nil
}

func PutGame(game *wordle.Game) {
	Games[game.Id] = game
}

func UpdateGame(game *wordle.Game) error {
	return nil
}

func ResetGames() {
	for k := range Games {
		delete(Games, k)
	}
}

func RemoveGame(gameID string) error {
	delete(Games, gameID)
	return nil
}
