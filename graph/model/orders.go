package model

import "time"

type BookingOrder struct {
	ID                  string `gorm:"type:char(36);primary_key"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time        `sql:"index"`
	BookingOrderDraftId string            `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	BookingOrderDraft   BookingOrderDraft `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
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
	Removed    bool
}
type MenuOrder struct {
	ID               string `gorm:"type:char(36);primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time     `sql:"index"`
	MenuOrderDraftId string         `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	MenuOrderDraft   MenuOrderDraft `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	// Booking draft orderer
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	// Important manage order status and udpate detail
	OrderStatus string
}
type MenuOrderDraft struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// booking draft manager
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`

	MenuId       string   `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	Menu         Menu     `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	MenuItemId   string   `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	MenuItem     MenuItem `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	InvitationId string   `gorm:"UNIQUE_INDEX:compositeindex;index;null"`

	Invitation Invitation `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	Complete   bool
	Removed    bool
}

type ProductPurchaseOrder struct {
	ID                          string `gorm:"type:char(36);primary_key"`
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	DeletedAt                   *time.Time                `sql:"index"`
	ProductPurchaseOrderDraftId string                    `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	ProductPurchaseOrderDraft   ProductPurchaseOrderDraft `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	// Booking draft orderer
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	// Important manage order status and udpate detail
	OrderStatus string
}
type ProductPurchaseOrderDraft struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// booking draft manager
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`

	ProductSalesId string       `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	ProductSales   ProductSales `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	Complete       bool
	Removed        bool
}
type EventTicketPurchaseOrder struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// booking draft manager
	UserId string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User   User   `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`

	EventID  string `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	Event    Events `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	Complete bool
	Removed  bool
}

// TODO figure out how to manage travels for friends in circle going to a spot with an order
// TODO manage payments wallets
// TODO manage Group Payments
// TODO build out cirlceChat
// TODO manage what happens after order status
