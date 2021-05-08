package model

import "time"

type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type Notifications struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
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

type NotificationChannel struct {
	ID        string `json:"id"`
	Name      string
	Observers map[string]struct {
		Username string
		Message  chan *string
	}
}
type InstantMessage struct {
	ID      string `json:"id"`
	Name    string
	Message Message
}
