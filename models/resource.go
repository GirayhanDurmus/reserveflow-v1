package models

import "gorm.io/gorm"

type Resource struct {
	gorm.Model

	Name        string `gorm:"size:120;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	Capacity    int    `gorm:"not null" json:"capacity"`
	IsActive    bool   `gorm:"not null;default:true" json:"is_active"`
}
