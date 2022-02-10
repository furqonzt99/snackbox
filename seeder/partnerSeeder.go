package seeder

import (
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

func PartnerSeeder(db *gorm.DB)  {
	password, _ := helper.Hashpwd("1234qwer")
	user := models.User{
		Name:         "User 1",
		Email:        "user1@gmail.com",
		Password:     password,
		Address:      "Jl Matraman No 13",
		City:         "Jakarta",
		Balance:      0,
		Role: "partner",
	}

	db.Create(&user)

	partner1 := models.Partner{
		UserID:        3,
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