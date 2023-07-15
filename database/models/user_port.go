package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPort struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PortName  string    `gorm:"uniqueIndex:user_port_key"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:user_port_key"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
