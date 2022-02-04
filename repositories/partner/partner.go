package partner

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type PartnerInterface interface {
	RequestPartner(partner models.Partner) (models.Partner, error)
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
