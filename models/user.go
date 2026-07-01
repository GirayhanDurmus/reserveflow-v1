package models

import "gorm.io/gorm"

const (
	UserRoleUser          = "user"
	UserRoleAdmin         = "super-admin"
	UserRoleResourceAdmin = "resource-admin"
	UserRoleManager       = "manager"
)

type User struct {
	gorm.Model
	FullName     string `gorm:"size:120;not null" json:"full_name"`
	Email        string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`

	RoleID uint `gorm:"index;not null" json:"role_id"`
	Role   Role `json:"role"`
}
