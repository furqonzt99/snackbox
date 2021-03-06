package models

import "gorm.io/gorm"

type Cashout struct {
	gorm.Model
	UserID uint
	IdempotenceKey string
	ExternalID string
	BankCode string
	AccountHolderName string
	AccountNumber string
	Amount float64
	Description string
	Status string
	User User
}