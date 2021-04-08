package model

import "time"

type BookingOrder struct {
	ID                string `gorm:"type:char(36);primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time        `sql:"index"`
	OrderDraftId      string            `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	BookingOrderDraft BookingOrderDraft `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	// Booking draft orderer
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	// Important manage order status and udpate detail
	OrderStatus string
}

type BookingOrderDraft struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// booking draft manager
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`

	BookingId      string       `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	Booking        Booking      `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	BookingOfferId string       `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	BookingOffer   BookingOffer `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	InvitationId   string       `gorm:"UNIQUE_INDEX:compositeindex;index;null"`

	Invitation Invitation `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	Complete   bool
}
type MenuOrder struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type ProductPurchaseOrder struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type EventTicketPurchaseOrder struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
