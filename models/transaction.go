package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uint
	PartnerID uint
	Buffet bool
	Quantity int
	DateTime time.Time
	Latitude float64
	Longtitude float64
	Distance float64
	TotalPrice float64
	InvoiceID string
	PaymentUrl string
	PaymentChannel string
	PaymentMethod string
	PaidAt time.Time
	Status string
	User User
	Partner Partner
	Products []Product `gorm:"many2many:detail_transactions;"`
}

