package rating

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type RatingInterface interface {
	Create(rating models.Rating) (models.Rating, error)
	Update(models.Rating) (models.Rating, error)
	IsCanGiveRating(userId, transactionId int) (models.Transaction, error)
	GetByTrxID(trxID int) (models.Rating, error) 
}

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (rr RatingRepository) IsCanGiveRating(userId, transactionId int) (models.Transaction, error) {
	var transaction models.Transaction

	const CONFIRM_STATUS = "CONFIRM"

	if err := rr.db.Where("user_id = ? AND transaction_id = ? AND status = ?", userId, transactionId, CONFIRM_STATUS).First(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (rr *RatingRepository) Create(rating models.Rating) (models.Rating, error) {

	if err := rr.db.Create(&rating).Error; err != nil {
		return rating, err
	}

	var r models.Rating

	rr.db.Preload("User").First(&r, "user_id = ? AND transaction_id = ?", &rating.UserID, &rating.TransactionID)

	return rating, nil
}

func (rr *RatingRepository) Update(rating models.Rating) (models.Rating, error) {
	var r models.Rating

	if err := rr.db.First(&r, "user_id = ? AND transaction_id = ?", rating.UserID, rating.TransactionID).Error; err != nil {
		return r, err
	}

	rr.db.Model(&r).Updates(rating)

	rr.db.Preload("User").First(&r, "user_id = ? AND transaction_id = ?", &rating.UserID, &rating.TransactionID)

	return r, nil
}

func (rr *RatingRepository) GetByTrxID(trxID int) (models.Rating, error) {
	var r models.Rating

	if err := rr.db.Preload("User").First(&r, "transaction_id = ?", trxID).Error; err != nil {
		return r, err
	}

	return r, nil
}
