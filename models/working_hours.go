package models

import "gorm.io/gorm"

const (
	DayMonday    = "monday"
	DayTuesday   = "tuesday"
	DayWednesday = "wednesday"
	DayThursday  = "thursday"
	DayFriday    = "friday"
	DaySaturday  = "saturday"
	DaySunday    = "sunday"
)

type WorkingHour struct {
	gorm.Model

	ResourceID uint     `gorm:"not null" json:"resource_id"`
	DayOfWeek  string   `gorm:"size:20;not null" json:"day_of_week"`
	OpenTime   string   `gorm:"size:5" json:"open_time"`
	CloseTime  string   `gorm:"size:5" json:"close_time"`
	IsClosed   bool     `gorm:"not null;default:false" json:"is_closed"`
	Resource   Resource `gorm:"foreignKey:ResourceID" json:"-"`
}
