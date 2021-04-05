package model

import (
	"time"
)

// type User struct {
// 	ID string `json:"id"`
// }

type User struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	SomeFlag  bool       `gorm:"column:some_flag;not null;default:true"`
	// main content goes here
	Email    string `gorm:"uniqueIndex,unique"`
	Password string

	// Profile  Profile
}

// Profile is the model for the profile table.
// type Profile struct {
// 	Base
// 	Name   string    `gorm:"column:name;size:128;not null;"`
// 	UserID uuid.UUID `gorm:"type:uuid;column:user_foreign_key;not null;"`
// }
