package bank

import (
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/disbursement"
	"gorm.io/gorm"
)

type BankInterface interface {
	GetAvailableBanks() ([]xendit.DisbursementBank, error)
}

type BankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) *BankRepository {
	return &BankRepository{db: db}
}

func (br *BankRepository) GetAvailableBanks() ([]xendit.DisbursementBank, error) {
	availableBanks, err := disbursement.GetAvailableBanks()
	if err != nil {
		return availableBanks, err
	}

	return availableBanks, nil
}
