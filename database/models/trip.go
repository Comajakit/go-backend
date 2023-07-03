package models

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripName  string
	TripCode  string    `gorm:"uniqueKey"`
	UserId    uuid.UUID `gorm:"type:uuid"`
	User      User      `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
