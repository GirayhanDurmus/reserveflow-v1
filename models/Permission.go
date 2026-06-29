package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Method   string `gorm:"size:10;not null" json:"method"`
	Endpoint string `gorm:"size:255;not null" json:"endpoint"`
}
