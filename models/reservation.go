package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ReservationStatusHeld      = "held"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusCancelled = "cancelled"
	ReservationStatusExpired   = "expired"
)

type Reservation struct {
	gorm.Model

	UserID     uint      `gorm:"not null" json:"user_id"`
	ResourceID uint      `gorm:"not null" json:"resource_id"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
	Status     string    `gorm:"size:30;not null;default:held" json:"status"`
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Resource Resource `gorm:"foreignKey:ResourceID" json:"resource"`
}
