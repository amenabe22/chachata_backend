package model

import "time"

type Booking struct {
	ID                string `gorm:"type:char(36);primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
	BookingName       string
	CoverPic          []BookingCoverPic `gorm:"polymorphic:Owner"`
	BookingOffers     []BookingOffer    `gorm:"polymorphic:Owner"`
	ServiceProviderId string            `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ServiceProvider   ServiceProvider   `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
}

type BookingCoverPic struct {
	// Do more details on bookings
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	BookingId string
	Tags      []BookingTag `gorm:"polymorphic:Owner"`
	OwnerID   string
	OwnerType string
}

type BookingTag struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	TagTitle  string
	OwnerID   string
	OwnerType string
}

type BookingOffer struct {
	ID                 string `gorm:"type:char(36);primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `sql:"index"`
	ServiceTitle       string
	ServiceDescription string
	Pricings           []Pricing          `gorm:"polymorphic:Owner"`
	Pics               []BookingOfferPics `gorm:"polymorphic:Owner"`
	OwnerID            string
	OwnerType          string
}

type BookingOfferPics struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	OwnerID   string
	OwnerType string
	Pic       string
}
type Pricing struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	OwnerID   string
	OwnerType string
	// Label is for the timestamp detail of the pricing
	Label          string
	Price          string
	ComparingPrice string
	Timed          bool
	TimedLabelId   string     `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	TimedLabel     TimedLabel `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
}

type TimedLabel struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	TypeTitle string
}
