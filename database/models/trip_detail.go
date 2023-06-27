package models

import (
	"time"
)

type TripDetail struct {
	ID               uint `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	TripCode         string       `gorm:"primaryKey"`
	Trip             Trip         `gorm:"foreignKey:TripCode"`
	TripCategoryCode string       `gorm:"primaryKey"`
	TripCategory     TripCategory `gorm:"foreignKey:TripCategoryCode"`
	TotalVacationDay uint         // Total number of vacation days for the trip
	TotalMember      uint         // Total number of members participating in the trip
	TotalCost        uint         // Total cost of the trip
	CostPerMember    uint         // Cost per member
	StartDate        time.Time    // Start date of the vacation period
	EndDate          time.Time    // End date of the vacation period
}
