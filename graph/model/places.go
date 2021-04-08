package model

import (
	"time"
)

type Place struct {
	ID                string `gorm:"type:char(36);primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
	PlaceTypeId       string     `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	PlaceType         PlaceType  `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	AdminOrigin       bool
	ServiceProviderId string          `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ServiceProvider   ServiceProvider `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	PlacePics         []PlacePic      `gorm:"polymorphic:Owner"`
}

type PlaceType struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	TypeName  string
	Image     string
}

type PlacePic struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Pic       string
	OwnerID   string
	OwnerType string
}
