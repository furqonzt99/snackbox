package partner

import (
	"errors"

	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type PartnerInterface interface {
	RequestPartner(partner models.Partner) (models.Partner, error)
	GetAllPartner() ([]models.Partner, error)
	DeletePartner(userId int) error
}

type PartnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepo(db *gorm.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

func (p *PartnerRepository) RequestPartner(partner models.Partner) (models.Partner, error) {
	err := p.db.Save(&partner).Error
	if err != nil {
		return partner, err
	}
	return partner, nil
}

func (p *PartnerRepository) GetAllPartner() ([]models.Partner, error) {
	var partner []models.Partner

	err := p.db.Find(&partner).Error
	if err != nil {
		return nil, err
	}

	return partner, nil
}

func (p *PartnerRepository) DeletePartner(userId int) error {
	var delete models.Partner

	err := p.db.Where("user_id = ?", userId).Delete(&delete).Error
	if err != nil {
		return errors.New("gak ketemu idnya")
	}
	return nil

}
