package models

import (
	"time"

	"github.com/google/uuid"
)

type TripDetail struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	TripID           uuid.UUID
	Trip             Trip `gorm:"foreignKey:TripID"`
	TripCategoryID   uint
	TripCategory     TripCategory `gorm:"foreignKey:TripCategoryID"`
	TotalVacationDay uint         // Total number of vacation days for the trip
	TotalMember      uint         // Total number of members participating in the trip
	TotalCost        uint         // Total cost of the trip
	CostPerMember    uint         // Cost per member
	StartDate        time.Time    // Start date of the vacation period
	EndDate          time.Time    // End date of the vacation period
}
