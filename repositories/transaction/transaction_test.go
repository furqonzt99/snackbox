package transaction_test

import (
	"log"
	"os"
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/furqonzt99/snackbox/repositories/product"
	"github.com/furqonzt99/snackbox/repositories/transaction"
	"github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	constants.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	constants.XENDIT_CALLBACK_TOKEN = os.Getenv("XENDIT_CALLBACK_TOKEN")

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
		Name:     "test",
		Photo:    "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Address:  "test",
		City:     "test",
		Balance:  0,
		Role:     "partner",
	}
	db.Create(&dummyUser)

	//CREATE USER2
	dummyUser2 := models.User{
		Name:     "test2",
		Photo:    "test2",
		Email:    "test2@gmail.com",
		Password: "test1234",
		Address:  "test2",
		City:     "test2",
		Balance:  0,
		Role:     "user",
	}
	db.Create(&dummyUser2)

	//CREATE PARTNER
	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Description:   "partner1",
		Latitude:      10,
		Longtitude:    10,
		Address:       "partner1",
		City:          "partner1",
		LegalDocument: "partner1",
		Status:        "active",
	}
	db.Create(&dummyPartner)

	//CREATE PRODUCT
	dummyProduct := models.Product{
		PartnerID:   1,
		Title:       "rendang",
		Image:       "iamge",
		Type:        "ricebox",
		Description: "enak",
		Price:       1000,
	}
	db.Create(&dummyProduct)

	dummyTransaction := models.Transaction{
		UserID:     1,
		PartnerID:  1,
		Buffet:     false,
		Quantity:   50,
		Latitude:   10,
		Longtitude: 10,
		Distance:   2,
		TotalPrice: 20000,
		Status:     "PENDING",
	}
	db.Create(&dummyTransaction)

	t.Run("create order", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 20000
		mockTransaction.InvoiceID = "suka"

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{1})
		assert.Equal(t, float64(20000), res.TotalPrice)
	})

	t.Run("create order invalid 1", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.ID = 1
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 30000

		mockTransaction.InvoiceID = "suka"

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{1})
		assert.Equal(t, float64(30000), res.TotalPrice)
	})

	t.Run("create order invalid 3", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 30000
		mockTransaction.InvoiceID = "suka"

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{0})
		assert.Equal(t, float64(30000), res.TotalPrice)
	})

	t.Run("create order", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 20000
		mockTransaction.InvoiceID = "suka"

		res, _ := transactionRepo.Order(mockTransaction, "test9@gmail.com", []int{1})
		assert.Equal(t, float64(20000), res.TotalPrice)
	})

	t.Run("create order success", func(t *testing.T) {
		var mockTransaction models.Transaction
		mockTransaction.UserID = 1
		mockTransaction.PartnerID = 1
		mockTransaction.TotalPrice = 20000
		mockTransaction.InvoiceID = "suka"
		mockTransaction.Quantity = 5

		res, _ := transactionRepo.Order(mockTransaction, "test2@gmail.com", []int{1})
		assert.Equal(t, float64(55000), res.TotalPrice)
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

func TestGetOneForUser(t *testing.T) {
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

	t.Run("test TestGetOneForUser", func(t *testing.T) {

		res, _ := transactionRepo.GetOneForUser(1, 2)
		assert.Equal(t, "PAID", res.Status)
	})
	t.Run("test TestGetOneForUser invalid", func(t *testing.T) {
		db.Migrator().DropTable(&models.Transaction{})
		res, _ := transactionRepo.GetOneForUser(1, 2)
		assert.Equal(t, float64(0), res.User.Balance)
	})
}

func TestGetPartnerFromProduct(t *testing.T) {
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

	t.Run("test GetPartnerFromProduct", func(t *testing.T) {

		res, _ := transactionRepo.GetPartnerFromProduct(1)
		assert.Equal(t, "partner1", res.BussinessName)
	})

	t.Run("test GetPartnerFromProduct invalid", func(t *testing.T) {
		db.Migrator().DropTable(&models.Partner{})
		res, _ := transactionRepo.GetPartnerFromProduct(1)
		assert.Equal(t, "", res.BussinessName)
	})

	t.Run("test GetPartnerFromProduct invalid", func(t *testing.T) {
		db.Migrator().DropTable(&models.Product{})
		res, _ := transactionRepo.GetPartnerFromProduct(1)
		assert.Equal(t, "", res.BussinessName)
	})
}

func TestGetOneForPartner(t *testing.T) {
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

	t.Run("test GetOneForPartner", func(t *testing.T) {

		res, _ := transactionRepo.GetOneForPartner(1, 1)
		assert.Equal(t, "PAID", res.Status)
	})
	t.Run("test GetOneForPartner invalid", func(t *testing.T) {
		db.Migrator().DropTable(&models.Transaction{})
		res, _ := transactionRepo.GetOneForPartner(1, 1)
		assert.Equal(t, uint(0), res.ID)
	})
}

func TestCallback(t *testing.T) {
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
		InvoiceID:  "11",
	}
	db.Create(&dummyTransaction)

	dummyTransaction2 := models.Transaction{
		PartnerID:  1,
		UserID:     99,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
		InvoiceID:  "22",
	}
	db.Create(&dummyTransaction2)

	dummyTransaction3 := models.Transaction{
		PartnerID:  1,
		UserID:     8585,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
		InvoiceID:  "23abc",
	}
	db.Create(&dummyTransaction3)

	t.Run("test Callback success", func(t *testing.T) {
		transaction := models.Transaction{}
		transaction.UserID = 2
		res, _ := transactionRepo.Callback("11", transaction, 500)
		assert.Equal(t, "", res.Partner.BussinessName)
	})
	t.Run("test Callback transaction not found", func(t *testing.T) {
		transaction := models.Transaction{}
		transaction.UserID = 2
		res, _ := transactionRepo.Callback("99", transaction, 500)
		assert.Equal(t, "", res.Partner.BussinessName)
	})

	t.Run("test Callback user not found", func(t *testing.T) { //////////////////////////<<<<<<
		transaction := models.Transaction{}
		transaction.UserID = 2
		res, _ := transactionRepo.Callback("23abc", dummyTransaction3, 500)
		assert.Equal(t, "", res.Partner.BussinessName)
	})

	t.Run("test Callback", func(t *testing.T) {
		transaction := models.Transaction{}
		transaction.UserID = 9
		res, _ := transactionRepo.Callback("11", transaction, 500)
		assert.Equal(t, "", res.Partner.BussinessName)
	})

}

func TestGetDistance(t *testing.T) {
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

	t.Run("test GetDistance success", func(t *testing.T) {

		res, _ := transactionRepo.GetDistance(1, 100, 100)
		assert.NotNil(t, res)
	})

	t.Run("test GetDistance partner invalid", func(t *testing.T) {

		res, _ := transactionRepo.GetDistance(99, 100, 100)
		assert.Equal(t, float64(0), res)
	})

}
