package models

import (
	"time"

	"github.com/google/uuid"
)

type ThemePercentagePair struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	StrategyProfileID uuid.UUID `gorm:"type:uuid"`
	Theme             string    `gorm:"type:varchar(255)"`
	Percentage        float64
}

type PortStrategyProfile struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID      uuid.UUID `gorm:"type:uuid"`
	User         User      `gorm:"foreignKey:OwnerID"`
	UserPortID   uuid.UUID `gorm:"type:uuid;uniqueIndex:port_strategy_key"`
	UserPort     UserPort  `gorm:"foreignKey:UserPortID"`
	StrategyName string
	Themes       []ThemePercentagePair `gorm:"foreignKey:StrategyProfileID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
