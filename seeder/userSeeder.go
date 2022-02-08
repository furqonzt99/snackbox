package seeder

import (
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB)  {
	password, _ := helper.Hashpwd("1234qwer")
	user1 := models.User{
		Name:         "User 1",
		Email:        "user1@gmail.com",
		Password:     password,
		Address:      "Jl Matraman No 13",
		City:         "Jakarta",
		Balance:      0,
	}

	db.Create(&user1)
}