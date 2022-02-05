package partner

import (
	"errors"

	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type PartnerInterface interface {
	RequestPartner(partner models.Partner) (models.Partner, error)
	GetAllPartner() ([]models.Partner, error)
	FindPartnerId(id int) (models.Partner, error)
	AcceptPartner(partner models.Partner) error
	RejectPartner(partner models.Partner) error
	GetAllPartnerProduct() ([]models.Partner, error)
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

func (p *PartnerRepository) FindPartnerId(id int) (models.Partner, error) {
	var fond models.Partner

	err := p.db.Where("id = ?", id).First(&fond).Error
	if err != nil {
		return fond, err
	}
	return fond, nil
}

func (p *PartnerRepository) AcceptPartner(partner models.Partner) error {

	if partner.Status == "active" {
		return errors.New("status already active")
	}
	partner.Status = "active"
	err := p.db.Save(&partner).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PartnerRepository) RejectPartner(partner models.Partner) error {

	if partner.Status == "reject" {
		return errors.New("status already reject")
	}
	partner.Status = "reject"
	err := p.db.Save(&partner).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PartnerRepository) GetAllPartnerProduct() ([]models.Partner, error) {
	var partner []models.Partner

	err := p.db.Preload("products").Find(&partner).Error
	if err != nil {
		return nil, err
	}
	return partner, nil
}
