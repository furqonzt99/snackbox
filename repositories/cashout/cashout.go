package cashout

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type CashoutInterface interface {
	Cashout(cashout models.Cashout) (models.Cashout, error)
	History(userID int) ([]models.Cashout, error)
	CheckBalance(userID int) (models.User, error)
	CallbackSuccess(extID string, cashout models.Cashout) (models.Cashout, error)
	CallbackFailed(extID string, cashout models.Cashout) (models.Cashout, error)
}

type CashoutRepository struct {
	db *gorm.DB
}

func NewCashoutRepository(db *gorm.DB) *CashoutRepository {
	return &CashoutRepository{db: db}
}

func (cr *CashoutRepository) Cashout(cashout models.Cashout) (models.Cashout, error) {
	var user models.User

	var err error

	cr.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&cashout).Error; err != nil {
			return err
		}

		if err = tx.First(&user, cashout.UserID).Error; err != nil {
			return err
		}

		newBalance := user.Balance - cashout.Amount

		if err = tx.Model(&user).Update("balance", newBalance).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return cashout, err
	}

	return cashout, nil
}

func (cr *CashoutRepository) History(userID int) ([]models.Cashout, error) {
	var cashouts []models.Cashout

	if err := cr.db.Where("user_id = ?", userID).Find(&cashouts).Error; err != nil {
		return nil, err
	}

	return cashouts, nil
}

func (cr *CashoutRepository) CheckBalance(userID int) (models.User, error) {
	var user models.User

	if err := cr.db.First(&user, userID).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (cr *CashoutRepository) CallbackSuccess(extID string, cashout models.Cashout) (models.Cashout, error) {

	var cashoutDB models.Cashout

	if err := cr.db.First(&cashoutDB, "external_id = ?", extID).Error; err != nil {
		return cashout, err
	}

	if err := cr.db.Model(&cashoutDB).Updates(cashout).Error; err != nil {
		return cashout, err
	}

	return cashout, nil
}

func (cr *CashoutRepository) CallbackFailed(extID string, cashout models.Cashout) (models.Cashout, error) {

	var user models.User
	var cashoutDB models.Cashout

	var err error
	cr.db.Transaction(func(tx *gorm.DB) error {

		if err = tx.First(&cashoutDB, "external_id = ?", extID).Error; err != nil {
			return err
		}

		if err = tx.First(&user, cashout.UserID).Error; err != nil {
			return err
		}

		newBalance := user.Balance + cashout.Amount

		if err = tx.Model(&user).Update("balance", newBalance).Error; err != nil {
			return err
		}

		if err = tx.Model(&cashoutDB).Updates(cashout).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return cashout, err
	}

	return cashout, nil
}
