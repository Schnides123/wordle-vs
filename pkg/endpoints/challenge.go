package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Schnides123/wordle-vs/pkg/data"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
	"github.com/Schnides123/wordle-vs/pkg/ws"
)

func ChallengeLinkHandler(w http.ResponseWriter, r *http.Request) {
	optionsString := r.URL.Query().Get("gameOptions")
	var options wordle.GameOptions
	if len(optionsString) == 0 {
		options = wordle.DefaultOptions
	} else {
		err := json.Unmarshal([]byte(optionsString), &options)
		if err != nil {
			options = wordle.DefaultOptions
		}
	}
	game := wordle.NewGame(options)
	data.PutGame(game)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// upgrader.Upgrade sets HTTP failure status code, so
		// just need to set response body
		fmt.Fprint(w, err.Error())
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	connection := ws.NewConnection(conn)
	connection.Start()

	id1 := generatePlayerID()
	client := ws.NewClient(id1, connection, game)
	client.Run()
}
