package transaction

import (
	"fmt"
	"strconv"

	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type TransactionInterface interface {
	Order(transaction models.Transaction, email string, products []int) (models.Transaction, error)
	Accept(trxID, partnerID int) (models.Transaction, error)
	Reject(trxID, partnerID int) (models.Transaction, error)
	Send(trxID, partnerID int) (models.Transaction, error)
	Confirm(trxID, userID int) (models.Transaction, error)
	GetAllForPartner(partnerID int) ([]models.Transaction, error)
	GetAllForUser(userID int) ([]models.Transaction, error)
	GetOneForUser(trxID, userID int) (models.Transaction, error)
	GetOneForPartner(trxID, partnerID int) (models.Transaction, error)
	GetDistance(partnerID int, latitude, longtitude float64) (float64, error)

	GetPartnerFromProduct(productID int) (models.Partner, error)
	Callback(invId string, transaction models.Transaction, refund float64) (models.Transaction, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) Order(transaction models.Transaction, email string, products []int) (models.Transaction, error) {
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		for _, product := range products {
			if err := tx.Create(&models.DetailTransaction{
				TransactionID: transaction.ID,
				ProductID:     uint(product),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return transaction, err
	}

	var user models.User
	if err := tr.db.First(&user, "email = ?", email).Error; err != nil {
		return transaction, err
	}

	err = tr.db.Transaction(func(tx *gorm.DB) error {

		if err := tr.db.Preload("Products").First(&transaction, transaction.ID).Error; err != nil {
			return err
		}

		var transactionPayment models.Transaction

		transactionPayment, err = helper.CreateInvoice(transaction, email, user.Balance)
		if err != nil {
			return err
		}

		balanceRemain := user.Balance - transactionPayment.TotalPrice
		if balanceRemain <= 0 {
			balanceRemain = 0
		}

		// update user balance
		if err := tx.Model(&user).Update("balance", balanceRemain).Error; err != nil {
			return err
		}

		if err := tx.Model(&transaction).Updates(transactionPayment).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return transaction, err
	}

	if err := tr.db.Preload("User").Preload("Products").First(&transaction, transaction.ID).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) Accept(trxID int, partnerID int) (models.Transaction, error) {
	trx := models.Transaction{}

	const PAID_STATUS = "PAID"
	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	const ACCEPT_STATUS = "ACCEPT"
	tr.db.Model(&trx).Update("status", ACCEPT_STATUS)

	return trx, nil
}

func (tr *TransactionRepository) Reject(trxID int, partnerID int) (models.Transaction, error) {
	trx := models.Transaction{}

	const PAID_STATUS = "PAID"
	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	tr.db.Transaction(func(tx *gorm.DB) error {

		const REJECT_STATUS = "REJECT"
		if err := tx.Model(&trx).Update("status", REJECT_STATUS).Error; err != nil {
			return err
		}

		user := models.User{}
		if err := tx.First(&user, trx.UserID).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Update("balance", trx.TotalPrice).Error; err != nil {
			return err
		}

		return nil
	})

	return trx, nil
}

func (tr *TransactionRepository) Send(trxID int, partnerID int) (models.Transaction, error) {
	trx := models.Transaction{}

	const ACCEPT_STATUS = "ACCEPT"
	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, ACCEPT_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	const SEND_STATUS = "SEND"
	tr.db.Model(&trx).Update("status", SEND_STATUS)

	return trx, nil
}

func (tr *TransactionRepository) Confirm(trxID int, userID int) (models.Transaction, error) {
	trx := models.Transaction{}

	const SEND_STATUS = "SEND"
	if err := tr.db.Where("user_id = ? AND status = ?", userID, SEND_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	tr.db.Transaction(func(tx *gorm.DB) error {

		const CONFIRM_STATUS = "CONFIRM"
		if err := tx.First(&trx, trxID).Update("status", CONFIRM_STATUS).Error; err != nil {
			return err
		}

		partner := models.Partner{}
		if err := tx.First(&partner, trx.PartnerID).Error; err != nil {
			return err
		}

		user := models.User{}
		if err := tx.First(&user, "id = ?", partner.UserID).Model(&user).Update("balance", trx.TotalPrice).Error; err != nil {
			return err
		}

		return nil
	})

	return trx, nil
}

func (tr *TransactionRepository) GetAllForPartner(partnerID int) ([]models.Transaction, error) {
	trx := []models.Transaction{}

	const PAID_STATUS = "PAID"

	if err := tr.db.Preload("User").Preload("Products").Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).Find(&trx).Error; err != nil {
		return nil, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetAllForUser(userID int) ([]models.Transaction, error) {
	trx := []models.Transaction{}

	if err := tr.db.Preload("User").Preload("Products").Where("user_id = ?", userID).Find(&trx).Error; err != nil {
		return nil, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetOneForUser(trxID, userID int) (models.Transaction, error) {
	trx := models.Transaction{}

	if err := tr.db.Preload("User").Preload("Products").Where("user_id = ?", userID).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetOneForPartner(trxID, partnerID int) (models.Transaction, error) {
	trx := models.Transaction{}

	const PAID_STATUS = "PAID"

	if err := tr.db.Preload("User").Preload("Products").Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetPartnerFromProduct(productID int) (models.Partner, error) {
	product := models.Product{}

	if err := tr.db.First(&product, productID).Error; err != nil {
		return models.Partner{}, err
	}

	partner := models.Partner{}

	if err := tr.db.First(&partner, product.PartnerID).Error; err != nil {
		return partner, err
	}

	return partner, nil
}

func (tr *TransactionRepository) Callback(invId string, transaction models.Transaction, refund float64) (models.Transaction, error) {

	var trx models.Transaction
	var user models.User

	if err := tr.db.First(&trx, "invoice_id = ?", invId).Error; err != nil {
		return transaction, err
	}

	if err := tr.db.First(&user, trx.UserID).Error; err != nil {
		return transaction, err
	}

	balance := user.Balance + refund

	err := tr.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&user).Update("balance", balance).Error; err != nil {
			return err
		}

		if err := tx.Model(&trx).Updates(transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) GetDistance(partnerID int, latitude, longtitude float64) (float64, error) {
	partner := models.Partner{}

	if err := tr.db.First(&partner, partnerID).Error; err != nil {
		return 0, err
	}

	var distance float64
	const EARTH_RADIUS_IN_KILOMETER = 6371

	if err := tr.db.Raw("select (? * acos ( cos ( radians( ? ) ) * cos ( radians (latitude) ) * cos ( radians (longtitude) - radians( ? ) ) + sin(radians( ? )) * sin(radians(latitude)))) as distance from partners where id = ?", EARTH_RADIUS_IN_KILOMETER, latitude, longtitude, latitude, partnerID).Scan(&distance).Error; err != nil {
		return 0, err
	}

	s := fmt.Sprintf("%.2f", distance)

	formatDistance, _ := strconv.ParseFloat(s, 64)

	return formatDistance, nil
}
