package product_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/furqonzt99/snackbox/repositories/product"
	"github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *user.UserRepository
var partnerRepo *partner.PartnerRepository
var productRepo *product.ProductRepository

func TestAddProduct(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = user.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)
	productRepo = product.NewProductRepo(db)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Partner{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Rating{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.DetailTransaction{})
	db.AutoMigrate(&models.Cashout{})

	//CREATE USER
	dummyUser := models.User{

		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE PRODUCT
	dummyProduct := models.Product{
		PartnerID:   1,
		Title:       "rendang",
		Type:        "ricebox",
		Description: "enak",
		Price:       1000,
	}
	productRepo.AddProduct(dummyProduct)

	t.Run("add product success", func(t *testing.T) {

		product := models.Product{
			Title: "jagung",
			Type:  "snack",
			Price: 1000,
		}

		res, _ := productRepo.AddProduct(product)
		assert.Equal(t, "jagung", res.Title)

	})

	t.Run("add product success", func(t *testing.T) {
		db.Migrator().DropTable(&models.Product{})
		product := models.Product{
			Title: "jagung",
			Type:  "snack",
			Price: 1000,
		}

		_, err := productRepo.AddProduct(product)
		assert.NotNil(t, err)

	})

}

func TestFindProduct(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = user.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)
	productRepo = product.NewProductRepo(db)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Partner{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Rating{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.DetailTransaction{})
	db.AutoMigrate(&models.Cashout{})

	//CREATE USER
	dummyUser := models.User{

		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE PRODUCT
	dummyProduct := models.Product{
		PartnerID:   1,
		Title:       "rendang",
		Type:        "ricebox",
		Description: "enak",
		Price:       1000,
	}
	productRepo.AddProduct(dummyProduct)

	t.Run("find product success", func(t *testing.T) {

		product := models.Product{
			PartnerID: 1,
			Title:     "jagung",
			Type:      "snack",
			Price:     1000,
		}
		productRepo.AddProduct(product)
		res, _ := productRepo.FindProduct(2, 1)
		assert.Equal(t, "jagung", res.Title)

	})

	t.Run("find product failed", func(t *testing.T) {

		product := models.Product{
			PartnerID: 1,
			Title:     "jagung",
			Type:      "snack",
			Price:     1000,
		}
		productRepo.AddProduct(product)
		res, _ := productRepo.FindProduct(3, 10)
		assert.Equal(t, "", res.Title)

	})

}

func TestDeleteProduct(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = user.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)
	productRepo = product.NewProductRepo(db)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Partner{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Rating{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.DetailTransaction{})
	db.AutoMigrate(&models.Cashout{})

	//CREATE USER
	dummyUser := models.User{

		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE PRODUCT
	dummyProduct := models.Product{
		PartnerID:   1,
		Title:       "rendang",
		Type:        "ricebox",
		Description: "enak",
		Price:       1000,
	}
	productRepo.AddProduct(dummyProduct)

	t.Run("find product success", func(t *testing.T) {

		product := models.Product{
			PartnerID: 1,
			Title:     "jagung",
			Type:      "snack",
			Price:     1000,
		}
		productRepo.AddProduct(product)
		err := productRepo.DeleteProduct(2, 1)
		assert.Nil(t, err)

	})

	t.Run("find product failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Product{})

		err := productRepo.DeleteProduct(300, 100)
		assert.NotNil(t, err)
	})

}

func TestGetAllProduct(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = user.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)
	productRepo = product.NewProductRepo(db)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Partner{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Rating{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.DetailTransaction{})
	db.AutoMigrate(&models.Cashout{})

	//CREATE USER
	dummyUser := models.User{

		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE PRODUCT
	dummyProduct := models.Product{
		PartnerID:   1,
		Title:       "rendang",
		Type:        "ricebox",
		Description: "enak",
		Price:       1000,
	}
	productRepo.AddProduct(dummyProduct)

	t.Run("get all product success", func(t *testing.T) {

		dummyProduct := models.Product{
			PartnerID:   1,
			Title:       "rendang",
			Type:        "ricebox",
			Description: "enak",
			Price:       1000,
		}
		productRepo.AddProduct(dummyProduct)
		res, _ := productRepo.GetAllProduct(1, 10, "rendang", "", 0, 0)
		assert.Equal(t, "rendang", res[0].Title)

	})
	t.Run("get all product failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Product{})

		res, _ := productRepo.GetAllProduct(11, 10, "nasi", "", 0, 0)
		assert.Equal(t, []models.Product(nil), res)

	})

}
func TestUploadImage(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = user.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)
	productRepo = product.NewProductRepo(db)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Partner{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Rating{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.DetailTransaction{})
	db.AutoMigrate(&models.Cashout{})

	//CREATE USER
	dummyUser := models.User{

		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE PRODUCT
	dummyProduct := models.Product{
		PartnerID:   1,
		Title:       "rendang",
		Type:        "ricebox",
		Description: "enak",
		Price:       1000,
	}
	productRepo.AddProduct(dummyProduct)

	t.Run("upload product iamge", func(t *testing.T) {

		product := models.Product{
			PartnerID: 1,
			Title:     "jagung",
			Type:      "snack",
			Price:     1000,
		}
		productRepo.AddProduct(product)
		res, _ := productRepo.UploadImage(1, product)
		assert.Equal(t, "jagung", res.Title)
	})

	t.Run("upload product iamge failed 1", func(t *testing.T) {

		product := models.Product{
			PartnerID: 1,
			Title:     "jagung",
			Type:      "snack",
			Price:     1000,
		}
		productRepo.AddProduct(product)
		res, _ := productRepo.UploadImage(4, product)
		assert.Equal(t, "", res.Title)
	})

	t.Run("upload product iamge failed 2", func(t *testing.T) {

		product := models.Product{
			PartnerID: 100,
			Title:     "jagung",
			Type:      "snack",
			Price:     1000,
		}
		productRepo.AddProduct(product)
		res, _ := productRepo.UploadImage(1, product)
		assert.Equal(t, "jagung", res.Title)
	})
}
