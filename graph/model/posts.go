package model

import (
	"time"
)

// main posts model
type Post struct {
	// TODO Add place Foreign key model here
	ID                string `gorm:"type:char(36);primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
	Caption           string
	UserId            string       `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	User              User         `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	PostImages        []PostImages `gorm:"polymorphic:Owner"`
	PostLikes         []Like       `gorm:"polymorphic:Owner"`
	Approved          bool
	IsDrafted         bool
	HashTags          []HashTag `gorm:"polymorphic:Owner"`
	Promotional       bool
	ServiceProviderId string          `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	ServiceProvider   ServiceProvider `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
}

// Comments section with a recursive setup
type Comment struct {
	ID            string `gorm:"type:char(36);primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Comment       string
	PostId        string `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Post          Post   `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	ParentComment string
}

// Posts images multiple relation with post
type PostImages struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	PostImage string
	OwnerID   string
	OwnerType string
}

type Like struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	UserId    string     `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	User      User       `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	OwnerID   string
	OwnerType string
}

type HashTag struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Tag       string
	OwnerID   string
	OwnerType string
}
