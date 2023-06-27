package models

import (
	"time"
)

type TripCategory struct {
	ID               uint `gorm:"primaryKey"`
	TripCategoryCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CategoryName     string
}
