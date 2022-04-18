package endpoints

import (
	"container/list"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Schnides123/wordle-vs/pkg/util"
	"github.com/Schnides123/wordle-vs/pkg/wordle"
	"github.com/Schnides123/wordle-vs/pkg/ws"
)

var (
	queue     = list.New()
	queueLock = sync.Mutex{}
	increment = 0
	once      = sync.Once{}
)

func MatchmakingHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure matchmaking thread is running
	once.Do(func() { go Matchmaker() })

	//upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("error upgrading connection")
		return
	}
	connection := ws.NewConnection(conn)
	connection.Start()

	//create channel for game
	channel := make(chan *ws.Event)
	defer connection.Close()

	//add to queue
	node := Enqueue(channel)

	//when dequeued, add connection to hub
	select {
	case err, _ := <-connection.ErrorChannel():
		// Error with connection. Disconnect.
		fmt.Printf("got connection err %v\n", err)
		// If there is a race  here, it should be treated as if
		// one of the paired users just left immediately from the game.
		RemoveFromQueue(node)
		return

	case <-connection.MessageChannel():
		// Client side erroneously sent a message before joining. This is
		// against protocol. Disconnect.
		RemoveFromQueue(node)
		return

	case event := <-channel:
		// We were finally put into a match. Join the hub
		if event.Type == ws.UpdateEvent {
			close(channel)
			client := ws.NewClient(event.PlayerID, connection, event.State)
			client.Run()

		} else if event.Type == ws.ErrorEvent {
			RemoveFromQueue(node)
			return
		}
	}
}

func RemoveFromQueue(el *list.Element) {
	queueLock.Lock()
	defer queueLock.Unlock()
	ch := queue.Remove(el).(chan *ws.Event)
	close(ch)
}

func Enqueue(channel chan *ws.Event) *list.Element {
	queueLock.Lock()
	defer queueLock.Unlock()
	return queue.PushBack(channel)
}

func Matchmaker() {
	for {
		//wait 50ms
		<-time.After(100 * time.Millisecond)

		//pull n connections from queue
		queueLock.Lock()

		for queue.Len() >= 2 {
			p1 := queue.Remove(queue.Front()).(chan *ws.Event)
			p2 := queue.Remove(queue.Front()).(chan *ws.Event)

			//assign playerIDs
			id1 := generatePlayerID()
			id2 := generatePlayerID()
			//create games
			game := wordle.NewGame(wordle.DefaultOptions)
			event1 := ws.Event{Type: ws.UpdateEvent, State: game, PlayerID: id1}
			event2 := ws.Event{Type: ws.UpdateEvent, State: game, PlayerID: id2}
			//send update event to connections
			select {
			case p1 <- &event1:
				//nothing
			default:
				//also nothing
			}
			select {
			case p2 <- &event2:
				//nothing
			default:
				//also nothing
			}
		}

		queueLock.Unlock()
	}
}

func generatePlayerID() string {
	var id string
	for i := 0; ; i++ {
		if increment > 0 {
			id = util.GetRandomWord(0) + "-" + util.GetRandomWord(0) + "-" + fmt.Sprint(increment)
		} else {
			id = util.GetRandomWord(0) + "-" + util.GetRandomWord(0)
		}

		if !ws.HasClientWithID(id) {
			return id
		}

		if i > 10000000 {
			increment++
			i = 0
		}
	}
}
