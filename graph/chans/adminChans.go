package chans

import "github.com/amenabe22/chachata_backend/graph/model"

type NotificationChannel struct {
	RoomId    string
	Message   string
	Observers map[string]struct {
		Username string
		Message  chan *model.Notifications
	}
}

type RoomNotification struct {
	RoomId    string
	Observers map[string]struct {
		Username string
		Message  chan *model.Notifications
	}
}
