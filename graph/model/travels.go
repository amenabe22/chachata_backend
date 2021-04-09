package model

import "time"

type UserTravel struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// TODO complete User travel setup
	OrderTravel            bool
	BookingOrderId         string               `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	BookingOrder           BookingOrder         `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	MenuOrderId            string               `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	MenuOrder              MenuOrder            `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	ProductPurchaseOrderId string               `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	ProductPurchaseOrder   ProductPurchaseOrder `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	TravelStatus           string
}
