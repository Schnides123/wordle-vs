//go:build !js && !wasm

package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

type R = *GuessRequest
type W = *Event

type GameConnection struct {
	connection     *websocket.Conn
	messageChannel chan R
	errorChannel   chan *Error
	writeLock      sync.Mutex
}

func NewConnection(conn *websocket.Conn) *GameConnection {
	return &GameConnection{
		connection:     conn,
		writeLock:      sync.Mutex{},
		messageChannel: make(chan R),
		errorChannel:   make(chan *Error),
	}
}

func (c *GameConnection) Start() {
	go c.Run()
}

func (c *GameConnection) Run() {
	defer c.connection.Close()
	defer close(c.errorChannel)
	defer close(c.messageChannel)

	for {
		var msg R
		_, r, err := c.connection.NextReader()
		if err != nil {
			c.errorChannel <- &Error{
				Details: err.Error(),
			}
			return
		}
		err = json.NewDecoder(r).Decode(&msg)
		if err == io.EOF {
			// One value is expected in the message.
			err = io.ErrUnexpectedEOF
		}

		if err != nil {
			c.errorChannel <- &Error{
				Details: fmt.Errorf("failed to parse message: %w", err).Error(),
			}
			continue
		}

		c.messageChannel <- msg
	}
}

func (c *GameConnection) Close() {
	// Unblocks the call to NextReader and causes Run() to break out
	c.connection.Close()
}

func (c *GameConnection) MessageChannel() <-chan R {
	return c.messageChannel
}

func (c *GameConnection) ErrorChannel() <-chan *Error {
	return c.errorChannel
}

func (c *GameConnection) Write(msg W) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	return c.connection.WriteJSON(msg)
}
