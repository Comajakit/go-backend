package models

import (
	"time"

	"github.com/google/uuid"
)

type PortStock struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerId           uuid.UUID `gorm:"type:uuid"`
	User              User      `gorm:"foreignKey:OwnerId"`
	UserPortID        uuid.UUID `gorm:"type:uuid;uniqueIndex:port_stock_key"`
	UserPort          UserPort  `gorm:"foreignKey:UserPortID"`
	Total             float64
	DivPerShare       float64
	DivInPercent      float64
	ExpectedDivReturn float64
	PercentageInPort  float64
	DivPercentPort    float64
	StockSymbol       string `gorm:"uniqueIndex:port_stock_key"`
	Volume            uint
	AveragePrice      float64
	StockType         string

	CreatedAt time.Time
	UpdatedAt time.Time
}
