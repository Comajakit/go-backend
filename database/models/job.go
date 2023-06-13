package models

import (
	"time"
)

type UserJob struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	JobTitle       string
	JobDescription string
	UserID         uint `gorm:"primaryKey"`
	User           User `gorm:"foreignKey:UserID"`
}
