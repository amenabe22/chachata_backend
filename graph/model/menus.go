package model

import "time"

type Menu struct {
	ID         string `gorm:"type:char(36);primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Name       string
	MenuTypeId string      `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	MenuType   MenuType    `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	MenuItem   []MenuItem  `gorm:"polymorphic:Owner"`
	Beverages  []Beverages `gorm:"polymorphic:Owner"`
	// Sp added to menu directly
	ServiceProviderId string          `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ServiceProvider   ServiceProvider `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
}

type Beverages struct {
	ID           string `gorm:"type:char(36);primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	Name         string
	Type         string
	ColdBeverage bool
	FreshDrink   bool
	OwnerID      string
	OwnerType    string
}
type MenuType struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	TypeName  string
}

type MenuItem struct {
	ID             string `gorm:"type:char(36);primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	ItemTitle      string
	Price          int
	Image          string
	ItemCategoryId string `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	ItemCategory   Post   `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
	InstantlyMade  bool
	Description    string
	Ingredients    []Ingredients `gorm:"polymorphic:Owner"`
	ItemType       string
	OwnerID        string
	OwnerType      string
}

type ItemCategory struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	CatTitle  string
}

type Ingredients struct {
	ID             string `gorm:"type:char(36);primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	IngredientName string
	OwnerID        string
	OwnerType      string
}
