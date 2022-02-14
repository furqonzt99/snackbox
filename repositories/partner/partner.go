package partner

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type PartnerInterface interface {
	ApplyPartner(partner models.Partner) (models.Partner, error)
	GetAllPartner() ([]models.Partner, error)
	GetPartner(partnerId int) (models.Partner, error)
	FindPartnerId(partnerId int) (models.Partner, error)
	FindUserId(userId int) (models.Partner, error)
	AcceptPartner(partner models.Partner) error
	RejectPartner(partner models.Partner) error
	GetAllPartnerProduct() ([]models.Partner, error)
	UploadDocument(partnerID int, partner models.Partner) (models.Partner, error)
	Report(partnerId int) ([]models.Transaction, error)
}

type PartnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepo(db *gorm.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

func (p *PartnerRepository) ApplyPartner(partner models.Partner) (models.Partner, error) {
	err := p.db.Save(&partner).Error
	if err != nil {
		return partner, err
	}
	return partner, nil
}

func (p *PartnerRepository) UploadDocument(partnerID int, partner models.Partner) (models.Partner, error) {
	var partnerDB models.Partner

	if err := p.db.First(&partnerDB, partnerID).Error; err != nil {
		return partnerDB, err
	}

	if err := p.db.Model(&partnerDB).Updates(partner).Error; err != nil {
		return partnerDB, err
	}

	return partnerDB, nil
}

func (p *PartnerRepository) GetAllPartner() ([]models.Partner, error) {
	var partner []models.Partner

	err := p.db.Find(&partner).Error
	if err != nil {
		return nil, err
	}

	return partner, nil
}

func (p *PartnerRepository) FindPartnerId(partnerId int) (models.Partner, error) {
	var partner models.Partner

	err := p.db.First(&partner, partnerId).Error
	if err != nil {
		return partner, err
	}
	return partner, nil
}

func (p *PartnerRepository) FindUserId(userId int) (models.Partner, error) {
	var partner models.Partner

	err := p.db.First(&partner, "user_id = ?", userId).Error
	if err != nil {
		return partner, err
	}

	return partner, nil
}

func (p *PartnerRepository) AcceptPartner(partner models.Partner) error {

	var user models.User
	partner.Status = "active"
	var err error
	p.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Save(&partner).Error
		if err != nil {
			return err
		}
		tx.First(&user, "id = ?", partner.UserID)
		user.Role = "partner"
		tx.Save(&user)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PartnerRepository) RejectPartner(partner models.Partner) error {

	partner.Status = "reject"
	err := p.db.Save(&partner).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *PartnerRepository) GetAllPartnerProduct() ([]models.Partner, error) {
	var partner []models.Partner

	err := p.db.Preload("Products").Find(&partner).Error
	if err != nil {
		return nil, err
	}
	return partner, nil
}

func (p *PartnerRepository) GetPartner(partnerId int) (models.Partner, error) {

	var partner models.Partner
	err := p.db.Preload("Products").First(&partner, partnerId).Error
	if err != nil {
		return partner, err
	}
	return partner, nil
}

func (p *PartnerRepository) Report(partnerId int) ([]models.Transaction, error) {

	var transaction []models.Transaction
	p.db.Order("created_at desc").Where("status <> ? AND status <> ?", "PENDING", "UNPAID").Preload("User").Preload("Partner").Preload("Products").Find(&transaction, "partner_id = ?", partnerId)

	return transaction, nil
}
