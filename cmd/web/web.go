//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"syscall/js"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/util"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
	"github.com/Schnides123/wordle-vs/pkg/ws"
)

var (
	APIBaseURL, _ = url.Parse("wss://warrdle.com")
)

func SlowJSValue(i interface{}) js.Value {
	// b, _ := json.Marshal(i)
	// var m map[string]interface{}
	// _ = json.Unmarshal(b, &m)
	// return js.ValueOf(m)

	v, _ := MashalJSRef(i)
	return v
}

func MashalJSRef(of interface{}) (js.Value, error) {
	switch of.(type) {
	case time.Duration:
		of = of.(time.Duration).String()
	case *util.Set:
		of = of.(*util.Set).Values()
	}

	v := reflect.ValueOf(of)
	t := v.Type()

	switch t.Kind() {
	case reflect.Bool:
		return js.ValueOf(of), nil
	case reflect.Int:
		return js.ValueOf(of), nil
	case reflect.Int8:
		return js.ValueOf(of), nil
	case reflect.Int16:
		return js.ValueOf(of), nil
	case reflect.Int32:
		return js.ValueOf(of), nil
	case reflect.Int64:
		return js.ValueOf(of), nil
	case reflect.Uint:
		return js.ValueOf(of), nil
	case reflect.Uint8:
		return js.ValueOf(of), nil
	case reflect.Uint16:
		return js.ValueOf(of), nil
	case reflect.Uint32:
		return js.ValueOf(of), nil
	case reflect.Uint64:
		return js.ValueOf(of), nil
	case reflect.Uintptr:
		return js.ValueOf(of), nil
	case reflect.Float32:
		return js.ValueOf(of), nil
	case reflect.Float64:
		return js.ValueOf(of), nil
	case reflect.Complex64:
		return js.ValueOf(of), nil
	case reflect.Complex128:
		return js.ValueOf(of), nil
	case reflect.Array:
		var err error
		// convert contents to jsValue
		ar := make([]interface{}, v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			ar[i], err = MashalJSRef(v.Index(i))
			if err != nil {
				return js.Null(), err
			}
		}
		return js.ValueOf(ar), nil
	case reflect.Chan:
		return js.Null(), errors.New("cannot marshal channel to jsref")
	case reflect.Func:
		if a, ok := of.(func(js.Value, []js.Value) interface{}); ok {
			return js.ValueOf(js.FuncOf(a)), nil
		}
		return js.Null(), errors.New("function has incorrect signature for conversion to js")
	case reflect.Interface:
		return js.Null(), errors.New("cannot marshal interface to jsref")
	case reflect.Map:
		res := js.ValueOf(map[string]interface{}{})

		for _, k := range v.MapKeys() {
			j, e := MashalJSRef(v.MapIndex(k))
			if e != nil {
				return js.Null(), e
			}
			res.Set(k.String(), j)
		}
		return res, nil
	case reflect.Ptr:
		return MashalJSRef(v.Elem().Interface())
	case reflect.Slice:
		res := js.Global().Get("Array").New(v.Len())
		for i := 0; i < v.Len(); i++ {
			m, err := MashalJSRef(v.Index(i).Interface())
			if err != nil {
				return js.Null(), err
			}
			res.SetIndex(i, m)
		}
		return res, nil
	case reflect.String:
		return js.ValueOf(of), nil
	case reflect.Struct:
		res := js.ValueOf(map[string]interface{}{})

		for i := 0; i < v.NumField(); i++ {
			fd := t.Field(i)
			f := v.Field(i)

			if fd.IsExported() {
				j, e := MashalJSRef(f.Interface())
				if e != nil {
					return js.Null(), e
				}
				res.Set(strings.ToLower(fd.Name), j)
			}
		}

		return res, nil
	case reflect.UnsafePointer:
		return js.Null(), errors.New("cannot marshal unsafe ptr to jsref")
	default:
		return js.Null(), errors.New("unrecognized reflect.Kind: " + t.Kind().String())
	}
}

type Session struct {
	conn   *TinyWebsocket
	state  *wordle.Game
	player string
}

func NewSession(conn *TinyWebsocket) *Session {
	return &Session{
		conn: conn,
	}
}

func (s *Session) Start() {
	go func() {
		defer s.conn.Close()
		println("Started session")

		decoder := json.NewDecoder(s.conn)
		for decoder.More() {
			println("start decode")
			var msg ws.Event
			if err := decoder.Decode(&msg); err != nil {
				if jsonErr, ok := err.(*json.SyntaxError); ok {
					problemPart := s.conn.allData[jsonErr.Offset-10 : jsonErr.Offset+10]
					println(string(s.conn.allData))
					err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
				}
				println("error parsing message as JSON:", err.Error())
				// println(string())
				return
			}

			switch msg.Type {
			case ws.ErrorEvent:
				js.Global().Get("OnGameError").Invoke(msg.Error.Details)
			case ws.UpdateEvent:
				if msg.PlayerID != s.player {
					s.player = msg.PlayerID
					js.Global().Get("OnPlayerID").Invoke(msg.PlayerID)
				}
				s.state = msg.State
				js.Global().Get("OnGameState").Invoke(SlowJSValue(msg.State))
			}
		}
		println("Stopped session")
	}()
}

func (s *Session) Guess(str string) {
	guess := ws.GuessRequest{
		Guess:  str,
		Player: s.player,
	}
	enc, e := json.Marshal(guess)
	if e != nil {
		js.Global().Get("OnGameError").Invoke(e.Error())
		return
	}
	_, e = s.conn.Write(enc)
	if e != nil {
		js.Global().Get("OnGameError").Invoke(e.Error())
		return
	}
}

var session *Session

func main() {
	js.Global().Set("NewGame", js.FuncOf(newGame))
	js.Global().Set("JoinGame", js.FuncOf(joinGame))
	js.Global().Set("JoinMatchmaking", js.FuncOf(joinMatchmaking))
	js.Global().Set("SubmitGuess", js.FuncOf(submitGuess))
	println("Game runtime loaded")

	// Wait forever
	select {}
}

func newGame(_ js.Value, _ []js.Value) interface{} {
	// start a new goroutine to stream objects from the server
	go func() {
		if session != nil {
			println("there's already an active session!")
			return
		}

		reqURL := *APIBaseURL
		reqURL.Path = "/blitz/challenge"
		w, err := NewTinyWebsocket(context.TODO(), reqURL.String())
		if err != nil {
			println(err)
			return
		}

		session = NewSession(w)
		session.Start()
	}()
	return nil
}

func joinGame(_ js.Value, args []js.Value) interface{} {
	// start a new goroutine to stream objects from the server
	go func() {
		if session != nil {
			if session.state != nil && session.state.Id == args[0].String() {
				println("there's already an active session!")
				return
			} else {
				session.conn.Close()
				session = nil
			}
		}

		reqURL := *APIBaseURL
		reqURL.Path = "/blitz/" + args[0].String()
		w, err := NewTinyWebsocket(context.TODO(), reqURL.String())
		if err != nil {
			println("error starting websocket", err)
			return
		}

		session = NewSession(w)
		session.Start()
	}()
	return nil
}

func joinMatchmaking(_ js.Value, _ []js.Value) interface{} {
	// start a new goroutine to stream objects from the server
	go func() {
		if session != nil {
			println("there's already an active session!")
			return
		}

		reqURL := *APIBaseURL
		reqURL.Path = "/blitz/"
		w, err := NewTinyWebsocket(context.TODO(), reqURL.String())
		if err != nil {
			println("rror starting websockt", err)
			return
		}

		session = NewSession(w)
		session.Start()
	}()
	return nil
}

func submitGuess(_ js.Value, args []js.Value) interface{} {
	// start a new goroutine to stream objects from the server
	go func() {
		if session == nil {
			println("no session to write guess")
			return
		}

		session.Guess(args[0].String())
	}()

	return nil
}
