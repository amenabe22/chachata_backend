package model

import "time"

type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type Chatroom struct {
	ID        string `json:"id"`
	Name      string
	Messages  []Message
	Observers map[string]struct {
		Username string
		Message  chan *Message
	}
}

type InstantMessage struct {
	ID      string `json:"id"`
	Name    string
	Message Message
}
