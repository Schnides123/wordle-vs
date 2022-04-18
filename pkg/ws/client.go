//go:build !js && !wasm

package ws

import (
	"fmt"
	"sync"

	"github.com/Schnides123/wordle-vs/pkg/data"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
)

// type GameConnection = Connection[*GuessRequest, *Event]

var (
	allClients = map[string]map[*Client]bool{}
	clientLock = sync.Mutex{}
)

func HasClientWithID(id string) bool {
	clientLock.Lock()
	defer clientLock.Unlock()
	if v, ok := allClients[id]; ok {
		return len(v) > 0
	}
	return false
}

type Client struct {
	id         string
	game       *wordle.Game
	connection *GameConnection
}

func NewClient(
	id string,
	connection *GameConnection,
	game *wordle.Game,
) *Client {
	res := &Client{
		id:         id,
		game:       game,
		connection: connection,
	}
	return res
}

func (c *Client) Update() {
	c.SendState()
}

func BroadcastGame(g *wordle.Game) {
	g.Cancer()
	vals := g.Players.Values()
	g.UnCancer()

	clientLock.Lock()
	defer clientLock.Unlock()

	for _, pId := range vals {
		if o, ok := allClients[pId]; ok {
			for v, _ := range o {
				v.Update()
			}
		}
	}
}

func (c *Client) Run() {
	func() {
		clientLock.Lock()
		defer clientLock.Unlock()
		if _, ok := allClients[c.id]; !ok {
			allClients[c.id] = make(map[*Client]bool)
		}
		allClients[c.id][c] = true
	}()

	defer func() {
		clientLock.Lock()
		defer clientLock.Unlock()
		delete(allClients[c.id], c)
	}()

	c.game.AddPlayer(c.id)
	BroadcastGame(c.game)

	defer func() {
		c.game.RemovePlayer(c.id)
		BroadcastGame(c.game)
	}()

	// Wait for guesses
	for {
		select {
		case msg, ok := <-c.connection.MessageChannel():
			if !ok {
				return
			}

			if len(msg.Guess) != c.game.Word.Length {
				c.SendError(fmt.Sprintf("guess lengh %v must match length of word %v", len(msg.Guess), c.game.Word.Length))
				continue
			}

			if err := c.game.Guess(msg.Guess, msg.Player); err != nil {
				c.SendError(err.Error())
				continue
			}

			data.UpdateGame(c.game)

			BroadcastGame(c.game)
		case err, ok := <-c.connection.ErrorChannel():
			if !ok {
				return
			}

			c.SendError(err.Error())
		}
	}
}

func (c *Client) SendError(error string) {
	c.SendEvent(&Event{
		Type:  ErrorEvent,
		Error: &Error{Details: error},
	})
}

func (c *Client) SendState() {
	c.SendEvent(&Event{
		Type:  UpdateEvent,
		State: c.game,
	})
}

func (c *Client) SendEvent(event *Event) {
	event.PlayerID = c.id

	c.game.Cancer()
	err := c.connection.Write(event)
	c.game.UnCancer()

	if err != nil {
		fmt.Printf("error sending message: %s\n", err.Error())
	}
}
