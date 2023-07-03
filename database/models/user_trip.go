package models

import (
	"time"

	"github.com/google/uuid"
)

type UserTrip struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripID    uuid.UUID `gorm:"type:uuid"`
	Trip      Trip      `gorm:"foreignKey:TripID"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
