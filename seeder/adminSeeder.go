package seeder

import (
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

func AdminSeeder(db *gorm.DB)  {
	password, _ := helper.Hashpwd("1234qwer")
	admin1 := models.User{
		Name:         "Admin 1",
		Email:        "admin1@snackbox.com",
		Password:     password,
		Address:      "Jl Matraman No 13",
		City:         "Jakarta",
		Balance:      0,
		Role:         "admin",
	}

	db.Create(&admin1)
}