package model

import "time"

type Events struct {
	ID                 string `gorm:"type:char(36);primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `sql:"index"`
	EventType          string
	EventName          string
	EventDescription   string
	EventPrice         int
	FreeEvent          bool
	EventPromotionPics []EventPromotionPic `gorm:"polymorphic:Owner"`
	MaxInvitations     int
	AddressId          string  `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Address            Address `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	Closed             bool
	EventDate          time.Time
}

type EventPromotionPic struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Pic       string
	OwnerID   string
	OwnerType string
}
