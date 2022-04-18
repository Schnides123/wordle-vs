package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/wordle"
	"github.com/Schnides123/wordle-vs/pkg/ws"
	"github.com/gorilla/websocket"
)

type User struct {
	playerID      string
	conn          *websocket.Conn
	state         wordle.Game
	updateChannel chan *ws.Event
	stopCh        chan bool
	wg            sync.WaitGroup
}

func NewUser(playerID string) *User {
	return &User{
		playerID:      playerID,
		updateChannel: make(chan *ws.Event),
		stopCh:        make(chan bool),
		wg:            sync.WaitGroup{},
	}
}

func (u *User) NewGame(options *wordle.GameOptions) error {
	reqUrl := url.URL{}
	reqUrl.Scheme = "ws"
	reqUrl.Host = serverURL.Host
	reqUrl.Path = "/blitz/challenge"

	query := reqUrl.Query()
	if options != nil {
		str, err := json.Marshal(options)
		if err != nil {
			return err
		}
		query.Add("gameOptions", string(str))
	}
	reqUrl.RawQuery = query.Encode()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(
		context.TODO(),
		reqUrl.String(),
		http.Header{},
	)

	if err != nil {
		return err
	}

	u.conn = conn

	// Start goroutine loop to read incoming messages
	u.run()

	return nil
}

func (u *User) JoinMatchmaking() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(
		context.TODO(),
		fmt.Sprintf("ws://%v/blitz/", serverURL.Host),
		http.Header{},
	)

	if err != nil {
		return err
	}

	u.conn = conn
	u.run()

	return nil
}

func (u *User) ConnectToGame(gameId string, options *wordle.GameOptions) error {
	reqUrl := url.URL{}
	reqUrl.Scheme = "ws"
	reqUrl.Host = serverURL.Host
	reqUrl.Path = "/blitz/" + gameId

	query := reqUrl.Query()
	if options != nil {
		str, err := json.Marshal(options)
		if err != nil {
			return err
		}
		query.Add("gameOptions", string(str))
	}
	query.Add("playerID", u.playerID)
	reqUrl.RawQuery = query.Encode()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(
		context.TODO(),
		reqUrl.String(),
		http.Header{},
	)

	if err != nil {
		return err
	}

	u.conn = conn

	// Start goroutine loop to read incoming messages
	u.run()

	return nil
}

func (u *User) run() {
	u.wg.Add(1)

	go func() {
		defer func() {
			u.wg.Done()
			u.Disconnect()
		}()

		for {
			var event ws.Event
			if err := u.conn.ReadJSON(&event); err != nil {
				return
			}

			select {
			case u.updateChannel <- &event:
				// Loop again
			case <-u.stopCh:
				return
			}
		}
	}()
}

func (u *User) Guess(word string) error {
	return u.conn.WriteJSON(ws.GuessRequest{
		Guess:  word,
		Player: u.playerID,
	})
}

func (u *User) Disconnect() error {
	if u.conn != nil {
		u.conn.WriteControl(
			websocket.CloseMessage,
			nil,
			time.Now().Add(200*time.Millisecond),
		)
		err := u.conn.Close()
		select {
		case u.stopCh <- true:
		default:
		}
		u.wg.Wait()
		return err
	}
	return nil
}

func (u *User) WaitForState(timeout time.Duration) error {
	select {
	case event := <-u.updateChannel:
		if event.Type == ws.UpdateEvent {
			u.state = *event.State
		}
		if event.Error != nil {
			return event.Error
		}
		if event.PlayerID != "" {
			u.playerID = event.PlayerID
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	}
}

func (u *User) State() *wordle.Game {
	return &u.state
}

func (u *User) GetPlayerID() string {
	return u.playerID
}

func (u *User) CanPlay() bool {
	return len(u.state.Guessed) == 0 ||
		u.state.Guessed[len(u.state.Guessed)-1].Player != u.playerID
}
