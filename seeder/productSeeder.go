package seeder

import (
	"fmt"

	"github.com/furqonzt99/snackbox/models"
	"gorm.io/gorm"
)

func ProductSeeder(db *gorm.DB)  {
	for i := 1; i <= 5; i++ {
		product := models.Product{
			PartnerID:   1,
			Title:       fmt.Sprint("Product ", i),
			Type:        "Rice Box",
			Description: "Descriptions",
			Price:       25000,
		}
		db.Create(&product)
	}
	
	for i := 6; i <= 10; i++ {
		product := models.Product{
			PartnerID:   1,
			Title:       fmt.Sprint("Product ", i),
			Type:        "Snack",
			Description: "Descriptions",
			Price:       2000,
		}
		db.Create(&product)
	}
}