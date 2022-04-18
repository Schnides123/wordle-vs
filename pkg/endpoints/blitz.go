package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Schnides123/wordle-vs/pkg/data"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
	"github.com/Schnides123/wordle-vs/pkg/ws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// For testing only. Orphans the old ws hub so that it is GC'd when
// all connections DC
func ResetBlitz() {
	//!TODO: remove all clients
}

func BlitzHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var gameID string
	var playerID string

	params := mux.Vars(r)
	gameID, _ = params["gameID"]

	if len(gameID) == 0 {
		//!TODO: validate GAMEID
		fmt.Printf("invalid gameID")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(w, "invalid gameID"); err != nil {
			fmt.Println("failed to write output")
		}
		return
	}

	if r.URL.Query().Has("playerID") {
		playerID = r.URL.Query().Get("playerID")
	} else {
		//!TODO check cookies
		playerID = generatePlayerID()
	}

	game, err := data.GetGame(gameID)
	if err != nil {
		// Write HTTP error response
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return

	} else if game == nil {
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
		game = wordle.NewGame(options)

		// force requested gameid for now
		game.Id = gameID

		data.PutGame(game)

		//!TODO: need a goroutine that cleans up games which have not had
		// any activity or are finished and not cleaned up
		//
		// Want to remove games with:
		//	0 players: should be handled when blit ws goes down, but should
		//				check anyway in case of non-gracefully handled dc.
		//  1 player: should we allow someone to wait indefinitely as long as
		//				WS is up?
		//  2 players: keep game running. but if both players DC i think we
		//				want a 5 minute timer on the game before it is wiped
	}

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

	client := ws.NewClient(playerID, connection, game)
	client.Run()
}
