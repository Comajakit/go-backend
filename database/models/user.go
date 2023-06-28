package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DisplayName string
	Username    string
	Password    string
	Email       string `gorm:"unique"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
