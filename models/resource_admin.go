package models

import "gorm.io/gorm"

type ResourceAdmin struct {
	gorm.Model

	UserID     uint `gorm:"not null;index" json:"user_id"`
	ResourceID uint `gorm:"not null;index" json:"resource_id"`

	User     User     `gorm:"foreignKey:UserID"     json:"user,omitempty"`
	Resource Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}
