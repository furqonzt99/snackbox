package models

import "gorm.io/gorm"

type Partner struct {
	gorm.Model
	UserID        uint `gorm:"unique"`
	BussinessName string
	Description   string
	Latitude      float64
	Longtitude    float64
	LegalDocument string
	Status        string `gorm:"default:pending"`
}
