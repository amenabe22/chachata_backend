package model

import "time"

type ProductSales struct {
	ID                 string `gorm:"type:char(36);primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time              `sql:"index"`
	ProductSalesTypeId string                  `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ProductSalesType   ProductSalesType        `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	Images             []ProductSalesTypeImage `gorm:"polymorphic:Owner"`
	SellingPrice       int
	ListingPrice       int
	CurrencyId         string   `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Currency           Currency `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	IsAvailable        bool
	Quantity           int
	ServiceProviderId  string          `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ServiceProvider    ServiceProvider `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
}

type Currency struct {
	ID           string `gorm:"type:char(36);primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	CurrencyName string
}
type ProductSalesType struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	TypeName  string
	Approved  bool
}
type ProductSalesTypeImage struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Image     string
	OwnerID   string
	OwnerType string
}
