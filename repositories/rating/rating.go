package rating

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type RatingInterface interface {
	Create(rating models.Rating) (models.Rating, error)
	Update(models.Rating) (models.Rating, error)
	IsCanGiveRating(userId, houseId int) (bool, error)
}

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (rr RatingRepository) IsCanGiveRating(userId, houseId int) (bool, error) {
	var transaction models.Transaction

	const PAID_STATUS = "PAID"

	if err := rr.db.Where("user_id = ? AND partner_id = ? AND status = ?", userId, houseId, PAID_STATUS).First(&transaction).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (rr *RatingRepository) Create(rating models.Rating) (models.Rating, error) {

	if err := rr.db.Create(&rating).Error; err != nil {
		return rating, err
	}

	var r models.Rating

	rr.db.Preload("User").First(&r, "user_id = ? AND partner_id = ?", &rating.UserID, &rating.PartnerID)

	return rating, nil
}

func (rr *RatingRepository) Update(rating models.Rating) (models.Rating, error) {
	var r models.Rating

	if err := rr.db.First(&r, "user_id = ? AND partner_id = ?", rating.UserID, rating.PartnerID).Error; err != nil {
		return r, err
	}

	rr.db.Model(&r).Updates(rating)

	rr.db.Preload("User").First(&r, "user_id = ? AND partner_id = ?", &rating.UserID, &rating.PartnerID)

	return r, nil
}
