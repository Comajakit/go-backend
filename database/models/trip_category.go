package models

import (
	"time"

	"github.com/google/uuid"
)

type TripCategory struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TripCategoryCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CategoryName     string
}
