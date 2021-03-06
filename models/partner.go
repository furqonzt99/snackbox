package models

import "gorm.io/gorm"

type Partner struct {
	gorm.Model
	UserID        uint `gorm:"primaryKey"`
	BussinessName string
	Description   string
	Latitude      float64
	Longtitude    float64
	Address       string
	City          string
	LegalDocument string
	Status        string `gorm:"default:DRAFT"`
	Products      []Product
	Ratings		  []Rating
}
