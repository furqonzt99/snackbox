package transaction_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/furqonzt99/snackbox/repositories/product"
	"github.com/furqonzt99/snackbox/repositories/transaction"
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
var transactionRepo *transaction.TransactionRepository

func TestOrder(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	t.Run("create order", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 20000

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{1})
		assert.Equal(t, float64(20000), res.TotalPrice)
	})

	t.Run("create order invalid 1", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.ID = 1
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 30000

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{1})
		assert.Equal(t, float64(30000), res.TotalPrice)
	})

	t.Run("create order invalid 3", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 30000

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{0})
		assert.Equal(t, float64(30000), res.TotalPrice)
	})

	t.Run("create order", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 20000

		res, _ := transactionRepo.Order(mockTransaction, "test9@gmail.com", []int{1})
		assert.Equal(t, float64(20000), res.TotalPrice)
	})

}

func TestAccept(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     2,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	t.Run("test accept", func(t *testing.T) {

		res, _ := transactionRepo.Accept(1, 1)
		assert.Equal(t, "ACCEPT", res.Status)
	})

	t.Run("test invalid", func(t *testing.T) {
		dummyTransaction := models.Transaction{
			PartnerID:  1,
			UserID:     2,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "UNPAID",
		}
		db.Create(&dummyTransaction)
		res, _ := transactionRepo.Accept(2, 1)
		assert.Equal(t, "", res.Status)
	})
}

func TestReject(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     2,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	t.Run("test reject", func(t *testing.T) {

		res, _ := transactionRepo.Reject(1, 1)
		assert.Equal(t, "REJECT", res.Status)
	})

	t.Run("test reject invalid 1", func(t *testing.T) {
		dummyTransaction := models.Transaction{
			PartnerID:  1,
			UserID:     2,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "UNPAID",
		}
		db.Create(&dummyTransaction)
		res, _ := transactionRepo.Reject(2, 1)
		assert.Equal(t, "", res.Status)
	})

	db.Create(&dummyTransaction)

	t.Run("test reject invalid 2", func(t *testing.T) {
		dummyTransaction2 := models.Transaction{
			PartnerID:  1,
			UserID:     2,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "ACCEPT",
		}
		db.Create(&dummyTransaction2)
		res, _ := transactionRepo.Reject(3, 1)
		assert.Equal(t, "", res.Status)
	})

}
func TestSend(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     2,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "ACCEPT",
	}
	db.Create(&dummyTransaction)

	t.Run("test send", func(t *testing.T) {

		res, _ := transactionRepo.Send(1, 1)
		assert.Equal(t, "SEND", res.Status)
	})

	t.Run("test send invalid", func(t *testing.T) {
		dummyTransaction2 := models.Transaction{
			PartnerID:  1,
			UserID:     2,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "UNPAID",
		}
		db.Create(&dummyTransaction2)
		res, _ := transactionRepo.Send(2, 1)
		assert.Equal(t, "", res.Status)
	})
}

func TestConfirm(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     2,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "SEND",
	}
	db.Create(&dummyTransaction)

	t.Run("test confirm", func(t *testing.T) {

		res, _ := transactionRepo.Confirm(1, 2)
		assert.Equal(t, "CONFIRM", res.Status)
	})

	t.Run("test confirm invalid 1", func(t *testing.T) {
		dummyTransaction2 := models.Transaction{
			PartnerID:  1,
			UserID:     2,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "SEND2",
		}
		db.Create(&dummyTransaction2)
		res, _ := transactionRepo.Confirm(2, 2)
		assert.Equal(t, "", res.Status)
	})

	t.Run("test confirm invalid 2", func(t *testing.T) { //salah
		dummyTransaction3 := models.Transaction{
			Model: gorm.Model{
				ID: 1,
			},
			PartnerID:  1,
			UserID:     3,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "SEND",
		}
		db.Create(&dummyTransaction3)
		res, _ := transactionRepo.Confirm(1, 3)
		assert.Equal(t, "", res.Status)
	})
}

func TestGetAllForPartner(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     2,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	t.Run("test GetAllForPartner", func(t *testing.T) {

		res, _ := transactionRepo.GetAllForPartner(1)
		assert.Equal(t, "PAID", res[0].Status)
	})
	t.Run("test GetAllForPartner", func(t *testing.T) {
		db.Migrator().DropTable(&models.Transaction{})
		res, _ := transactionRepo.GetAllForPartner(1)
		assert.Equal(t, []models.Transaction([]models.Transaction(nil)), res)
	})
}

func TestGetAllForUser(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)

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
		Role:     "partner",
	}
	userRepo.Register(dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{

		Email:    "test2@gmail.com",
		Password: "test1234",
		Role:     "user",
	}
	userRepo.Register(dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "active",
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

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     2,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	t.Run("test GetAllForUser", func(t *testing.T) {

		res, _ := transactionRepo.GetAllForUser(2)
		assert.Equal(t, "PAID", res[0].Status)
	})
	t.Run("test GetAllForUser", func(t *testing.T) {
		db.Migrator().DropTable(&models.Transaction{})
		res, _ := transactionRepo.GetAllForUser(1)
		assert.Equal(t, []models.Transaction([]models.Transaction(nil)), res)
	})
}
