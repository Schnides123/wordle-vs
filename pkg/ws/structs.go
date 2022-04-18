package ws

import "github.com/Schnides123/wordle-vs/pkg/wordle"

type Message struct {
	UID  string `json:"uid"`
	Data string `json:"data"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Board   Board  `json:"board"`
}

type EventType int

const (
	UpdateEvent EventType = iota
	ErrorEvent
)

type Event struct {
	Type     EventType    `json:"type"`
	State    *wordle.Game `json:"state,omitempty"`
	Error    *Error       `json:"error,omitempty"`
	PlayerID string       `json:"playerID,omitempty"`
}

type Error struct {
	Details string `json:"details"`
}

func (e *Error) Error() string {
	return e.Details
}

type Board struct {
	Word string `json:"word"`
	Rows []Row  `json:"rows"`
}

type Row struct {
	Letters string `json:"value"`
	Result  []int  `json:"result"`
}

type GuessRequest struct {
	Guess  string `json:"guess"`
	Player string `json:"player"`
}
