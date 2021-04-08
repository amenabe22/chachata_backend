package model

import "time"

type ServiceProvider struct {
	ID                  string `gorm:"type:char(36);primary_key"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time `sql:"index"`
	Phone               string
	SpID                string `gorm:"unique"`
	SpType              string
	Email               string `gorm:"unique"`
	BusinessName        string
	AddressId           string  `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Address             Address `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	BusinessPhone       string  `gorm:"unique"`
	BusinessDescription string
	IsVerified          bool
	AccountAppoved      bool
}

type Address struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// TODO make sure city and country matches
	City      string
	CountryId string  `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Country   Country `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
}

type Country struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string     `gorm:"unique"`
}

type PromotionalPosts struct {
	ID                string `gorm:"type:char(36);primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time      `sql:"index"`
	PostId            string          `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Post              Post            `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	ServiceProviderId string          `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ServiceProvider   ServiceProvider `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	Featured          bool
}
