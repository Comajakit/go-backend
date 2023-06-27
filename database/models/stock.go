package models

import (
	"time"
)

type Stock struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Code      string
	Email     string
}
