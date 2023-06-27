package models

import (
	"time"
)

type User struct {
	ID          uint `gorm:"primaryKey"`
	DisplayName string
	Username    string
	Password    string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
