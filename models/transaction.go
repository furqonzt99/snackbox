package models

import (
	"errors"
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
	Distance float64 `gorm:"default:null"`
	TotalPrice float64
	InvoiceID string
	PaymentUrl string
	PaymentChannel string
	PaymentMethod string
	PaidAt time.Time `gorm:"default:null"`
	Status string `gorm:"default:pending"`
	User User
	Partner Partner
	Products []Product `gorm:"many2many:detail_transactions;"`
}

type DetailTransaction struct {
  TransactionID  uint `gorm:"primaryKey"`
  ProductID uint `gorm:"primaryKey"`
}

func (DetailTransaction) BeforeCreate(db *gorm.DB) error {
  err := db.SetupJoinTable(&Transaction{}, "Products", &DetailTransaction{})

  if err != nil {
	  return errors.New(err.Error())
  }

  return nil
}
