package product

import (
	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

type ProductInterface interface {
	AddProduct(product models.Product) (models.Product, error)
	FindProduct(productId, partnerId int) (models.Product, error)
	DeleteProduct(productId, partnerId int) error
	GetAllProduct(offset, pageSize int, search, category string, latitude, longtitude float64) ([]models.Product, error)
	UploadImage(productID int, product models.Product) (models.Product, error)
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product models.Product) (models.Product, error) {

	err := p.db.Save(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (p *ProductRepository) UploadImage(productID int, product models.Product) (models.Product, error) {
	var productDB models.Product

	if err := p.db.First(&productDB, productID).Error; err != nil {
		return productDB, err
	}

	if err := p.db.Model(&productDB).Updates(product).Error; err != nil {
		return productDB, err
	}

	return productDB, nil
}

func (p *ProductRepository) FindProduct(productId, partnerId int) (models.Product, error) {
	var product models.Product
	err := p.db.Where("id = ? AND partner_id = ?", productId, partnerId).First(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (p *ProductRepository) DeleteProduct(productId, partnerId int) error {

	var delete models.Product
	err := p.db.Where("id = ? AND partner_id = ?", productId, partnerId).Delete(&delete).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) GetAllProduct(offset, pageSize int, search, category string, latitude, longtitude float64) ([]models.Product, error) {
	var products []models.Product

	nearestPartner := []int{}
	const EARTH_RADIUS_IN_KILOMETER = 6371 
	const MAX_DISTANCE = 10 

	if latitude != 0 && longtitude != 0 {
		p.db.Raw("select id, (? * acos ( cos ( radians( ? ) ) * cos ( radians (latitude) ) * cos ( radians (longtitude) - radians( ? ) ) + sin(radians( ? )) * sin(radians(latitude)))) as distance from partners having distance < ? order by distance", EARTH_RADIUS_IN_KILOMETER, latitude, longtitude, latitude, MAX_DISTANCE).Scan(&nearestPartner)
		p.db.Find(&products, "partner_id IN ?", nearestPartner)
	}

	err := p.db.Offset(offset).Limit(pageSize).Where("title LIKE ? AND type LIKE ?", "%"+search+"%", "%"+category+"%").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
