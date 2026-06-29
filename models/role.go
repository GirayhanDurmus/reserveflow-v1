package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string       `json:"description"`
	Permission  []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
