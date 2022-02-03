package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Email string
	Password string
	Address string
	City string
	Balance string
	Role string
	Partner Partner
	Transactions []Transaction
}