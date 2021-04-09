package model

import "time"

// OwnerID   string
// OwnerType string
// Pricings           []Pricing          `gorm:"polymorphic:Owner"`

// TimedLabelId   string          `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
// TimedLabel     ServiceProvider `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`

type Invitation struct {
	ID               string `gorm:"type:char(36);primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	CircleInvitation bool
	// add single circle for invitation
	CircleId string     `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	Circle   UserCircle `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
	UserId   string     `gorm:"UNIQUE_INDEX:compositeindex;index;null"`
	User     User       `gorm:"UNIQUE_INDEX:compositeindex;type:text;null"`
}
