package transaction

import (
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

	GetPartnerFromProduct(productID int) (models.Partner, error)
	Callback(invId string, transaction models.Transaction) (models.Transaction, error)
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
	
	err = tr.db.Transaction(func(tx *gorm.DB) error {

		if err := tr.db.Preload("Products").First(&transaction, transaction.ID).Error; err != nil {
			return err
		}

		transactionPayment, err := helper.CreateInvoice(transaction, email)
		if err != nil {
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

	return transaction, nil
}

func (tr *TransactionRepository) Accept(trxID int, partnerID int) (models.Transaction, error)  {
	trx := models.Transaction{}

	const PAID_STATUS = "PAID"
	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	const ACCEPT_STATUS = "accept"
	tr.db.Model(&trx).Update("status", ACCEPT_STATUS)

	return trx, nil
}

func (tr *TransactionRepository) Reject(trxID int, partnerID int) (models.Transaction, error)  {
	trx := models.Transaction{}

	const PAID_STATUS = "PAID"
	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	tr.db.Transaction(func(tx *gorm.DB) error {

		const REJECT_STATUS = "reject"
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

func (tr *TransactionRepository) Send(trxID int, partnerID int) (models.Transaction, error)  {
	trx := models.Transaction{}

	const ACCEPT_STATUS = "accept"
	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, ACCEPT_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	const SEND_STATUS = "send"
	tr.db.Model(&trx).Update("status", SEND_STATUS)

	return trx, nil
}

func (tr *TransactionRepository) Confirm(trxID int, userID int) (models.Transaction, error)  {
	trx := models.Transaction{}

	const SEND_STATUS = "send"
	if err := tr.db.Where("user_id = ? AND status = ?", userID, SEND_STATUS).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	tr.db.Transaction(func(tx *gorm.DB) error {

		const CONFIRM_STATUS = "confirm"
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

	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).Find(&trx).Error; err != nil {
		return nil, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetAllForUser(userID int) ([]models.Transaction, error) {
	trx := []models.Transaction{}

	if err := tr.db.Where("user_id = ?", userID).Find(&trx).Error; err != nil {
		return nil, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetOneForUser(trxID, userID int) (models.Transaction, error) {
	trx := models.Transaction{}

	if err := tr.db.Where("user_id = ?", userID).First(&trx, trxID).Error; err != nil {
		return trx, err
	}

	return trx, nil
}

func (tr *TransactionRepository) GetOneForPartner(trxID, partnerID int) (models.Transaction, error) {
	trx := models.Transaction{}

	const PAID_STATUS = "PAID"

	if err := tr.db.Where("partner_id = ? AND status = ?", partnerID, PAID_STATUS).First(&trx, trxID).Error; err != nil {
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

func (tr *TransactionRepository) Callback(invId string, transaction models.Transaction) (models.Transaction, error) {

	var trx models.Transaction

	if err := tr.db.First(&trx, "invoice_id = ?", invId).Error; err != nil {
		return transaction, err
	}

	if err := tr.db.Model(&trx).Updates(transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}