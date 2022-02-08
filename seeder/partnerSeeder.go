package seeder

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

func PartnerSeeder(db *gorm.DB)  {
	partner1 := models.Partner{
		UserID:        2,
		BussinessName: "Melati Katering",
		Description:   "Katering Murah Rasa Bintang 5",
		Latitude:      -7.7343187,
		Longtitude:    111.3404542,
		Address:       "Jl Matraman No 13",
		City:          "Jakarta",
		LegalDocument: "legal.pdf",
		Status:        "active",
	}

	db.Create(&partner1)
}