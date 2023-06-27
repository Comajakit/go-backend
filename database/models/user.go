package models

import (
	"time"
)

type User struct {
	ID          uint `gorm:"primaryKey"`
	DisplayName string
	Username    string
	Password    string
	Email       string `gorm:"unique"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
