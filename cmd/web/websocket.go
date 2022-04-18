//go:build js && wasm

// Must use this instead of other websocket libraries since TinyGo
// does not support time callbacks in its runtime
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"syscall/js"
)

type TinyWebsocket struct {
	ctx      context.Context
	cancel   context.CancelFunc
	socket   js.Value
	messages chan []byte
	url      string
	leftover []byte
	wg       sync.WaitGroup
	allData  []byte
}

var _ io.Reader = &TinyWebsocket{}
var _ io.Writer = &TinyWebsocket{}
var _ io.Closer = &TinyWebsocket{}

func NewTinyWebsocket(ctx context.Context, url string) (*TinyWebsocket, error) {
	fmt.Println("starting websocket:", url)
	ws := js.Global().Get("WebSocket")
	if ws.IsNull() {
		return nil, errors.New("could not get js global WebSocket")
	}
	ctx, cancel := context.WithCancel(ctx)

	res := &TinyWebsocket{
		socket:   ws,
		url:      url,
		messages: make(chan []byte, 10),
		leftover: []byte{},
		ctx:      ctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
		allData:  []byte{},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	onopen := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		fmt.Println("Connection successfully established")
		wg.Done()
		return nil
	})

	onerr := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		fmt.Println("error opening ws", args)
		wg.Done()
		return nil
	})

	res.socket = res.socket.New(res.url)
	res.socket.Call("addEventListener", "close", js.FuncOf(res.OnClose))
	res.socket.Call("addEventListener", "message", js.FuncOf(res.OnMessage))
	res.socket.Call("addEventListener", "error", onerr)
	res.socket.Call("addEventListener", "open", onopen)

	go func() {
		select {
		case <-ctx.Done():
			res.socket.Call("close", 1000)
			res.wg.Wait()
			// wait for all goroutines posting to messages to finish
			close(res.messages)
		}
	}()

	// Wait for open/error
	wg.Wait()

	// // Remove temporary error handler and replace with instance handler
	res.socket.Call("removeEventListener", "error", onerr)
	res.socket.Call("removeEventListener", "open", onopen)
	res.socket.Call("addEventListener", "error", js.FuncOf(res.OnError))

	return res, nil
}

func (w *TinyWebsocket) OnError(this js.Value, args []js.Value) interface{} {
	println("Connection has been unexpectedly terminated")
	return nil
}

func (w *TinyWebsocket) OnClose(this js.Value, args []js.Value) interface{} {
	w.Close()
	println("Connection has been closed")
	return nil
}

func (w *TinyWebsocket) OnMessage(this js.Value, args []js.Value) interface{} {
	println("Connection received message")
	// didSpawnGoRoutine := sync.WaitGroup{}
	// didSpawnGoRoutine.Add(1)
	w.wg.Add(1)

	// go func() {
	defer w.wg.Done()

	v := args[0].Get("data").String()
	// fmt.Println("received", v)

	// didSpawnGoRoutine.Done()
	select {
	case w.messages <- []byte(v):
		// yay sent
	case <-w.ctx.Done():
		// cancelled
	}
	// }()

	// didSpawnGoRoutine.Wait()
	return nil
}

func (w *TinyWebsocket) Close() error {
	w.cancel()
	return nil
}

func (w *TinyWebsocket) Read(b []byte) (nread int, err error) {
	defer func() {
		println("read", string(b))
		println("leftover", string(w.leftover))

		if err == nil {
			w.allData = append(w.allData, b[:nread]...)
		}
	}()

	if len(w.leftover) > 0 {
		nread = nread + copy(b, w.leftover)
		w.leftover = w.leftover[nread:]
		println("copied leftover", nread, len(w.leftover))
	}

loop:
	for len(b)-nread > 0 {
		select {
		case msg := <-w.messages:
			println("copying", nread, len(msg), len(b))
			copied := copy(b[nread:], msg)
			nread = nread + copied

			if copied < len(msg) {
				w.leftover = append(w.leftover, msg[copied:]...)
			}
			continue
		case <-w.ctx.Done():
			return nread, io.EOF
		default:
			break loop
		}
	}

	if len(b) > 0 && nread == 0 {
		if len(w.leftover) > 0 {
			nread = nread + copy(b[nread:], w.leftover)
			w.leftover = w.leftover[nread:]
			println("copied leftover", nread, len(w.leftover))
		}

		select {
		case msg := <-w.messages:
			println("copying", nread, len(msg), len(b))
			copied := copy(b[nread:], msg)
			nread = nread + copied

			if copied < len(msg) {
				w.leftover = append(w.leftover, msg[copied:]...)
			}
			break
		case <-w.ctx.Done():
			return nread, io.EOF
		}
	}
	return nread, nil
}

func (w *TinyWebsocket) Write(b []byte) (int, error) {
	dst := js.Global().Get("Uint8Array").New(len(b))
	copied := js.CopyBytesToJS(dst, b)
	w.socket.Call("send", dst)
	return copied, nil
}

func (w *TinyWebsocket) AllData() []byte {
	return w.allData
}
