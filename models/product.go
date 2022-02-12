package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	PartnerID   uint
	Title       string
	Image		string
	Type        string
	Description string
	Price       float64
	Partner     Partner
}
