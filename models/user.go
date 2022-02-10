package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique"`
	Password     string
	Address      string
	City         string
	Balance      float64 `gorm:"default:0"`
	Role         string  `gorm:"default:user"`
	Partner      Partner
	Transactions []Transaction
}
