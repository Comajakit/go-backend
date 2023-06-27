package models

import (
	"time"
)

type Trip struct {
	ID        uint   `gorm:"primaryKey"`
	TripCode  string `gorm:"uniqueKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
