package model

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	SomeFlag  bool       `gorm:"column:some_flag;not null;default:true"`
	// main content goes here
	Email      string `gorm:"uniqueIndex,unique"`
	Password   string
	ProfileId  string
	Profile    Profile
	IsVerified bool
	Qrcode     string
	// UserDevices []Devices `gorm:"many2many:devices;" json:"devices,omitempty"`
	UserDevices []Devices `gorm:"polymorphic:Owner"`
	OwnerID     string
	OwnerType   string
}

// Profile is the model for the profile table.
type Profile struct {
	ID         string `gorm:"type:char(36);primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Name       string     `gorm:"column:name;size:128;not null;"`
	Username   string     `gorm:"uniqueIndex,unique"`
	Phone      string     `gorm:"uniqueIndex,unique"`
	ProfilePic string
	Complete   bool
	Progress   int
	Bio        string
}

type Devices struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// UserID     string     `json:"user_id" gorm:"user_id"` //nolint:gofmt
	AppID      string
	DeviceName string
	OwnerID    string
	OwnerType  string
}
