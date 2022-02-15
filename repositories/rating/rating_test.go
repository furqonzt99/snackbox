package rating_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/furqonzt99/snackbox/repositories/product"
	"github.com/furqonzt99/snackbox/repositories/rating"
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
var ratingRepo *rating.RatingRepository

func TestCreate(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)
	ratingRepo = rating.NewRatingRepository(db)

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
		Status:        "active",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     1,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	t.Run("create rating", func(t *testing.T) {
		var mockRating models.Rating
		mockRating.TransactionID = 1
		mockRating.PartnerID = 1
		mockRating.UserID = 1
		mockRating.Rating = 5

		res, _ := ratingRepo.Create(mockRating)
		assert.Equal(t, 5, res.Rating)
	})

	t.Run("create rating failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Rating{})
		var mockRating models.Rating
		mockRating.TransactionID = 1
		mockRating.PartnerID = 1
		mockRating.UserID = 1
		mockRating.Rating = 5

		_, err := ratingRepo.Create(mockRating)
		// assert.Equal(t, 5, res.Rating)
		assert.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)
	ratingRepo = rating.NewRatingRepository(db)

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
		Status:        "active",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     1,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	dummyRating := models.Rating{
		TransactionID: 1,
		PartnerID:     1,
		UserID:        1,
		Rating:        5,
	}
	db.Create(&dummyRating)

	t.Run("update rating", func(t *testing.T) {
		var mockRating models.Rating
		mockRating.TransactionID = 1
		mockRating.PartnerID = 1
		mockRating.UserID = 1
		mockRating.Rating = 5

		res, _ := ratingRepo.Update(mockRating)
		assert.Equal(t, 5, res.Rating)
	})

	t.Run("update rating failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Rating{})
		var mockRating models.Rating
		mockRating.TransactionID = 99
		mockRating.PartnerID = 99
		mockRating.UserID = 1
		mockRating.Rating = 5

		_, err := ratingRepo.Update(mockRating)
		// assert.Equal(t, 5, res.Rating)
		assert.NotNil(t, err)
	})
}

func TestGetByTrxID(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)
	ratingRepo = rating.NewRatingRepository(db)

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
		Status:        "active",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     1,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	dummyRating := models.Rating{
		TransactionID: 1,
		PartnerID:     1,
		UserID:        1,
		Rating:        5,
	}
	db.Create(&dummyRating)

	t.Run("get rating Id", func(t *testing.T) {
		var mockRating models.Rating
		mockRating.TransactionID = 1
		mockRating.PartnerID = 1
		mockRating.UserID = 1
		mockRating.Rating = 5

		res, _ := ratingRepo.GetByTrxID(1)
		assert.Equal(t, 5, res.Rating)
	})

	t.Run("update rating failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Rating{})
		var mockRating models.Rating
		mockRating.TransactionID = 99
		mockRating.PartnerID = 99
		mockRating.UserID = 1
		mockRating.Rating = 5

		_, err := ratingRepo.GetByTrxID(1)
		// assert.Equal(t, 5, res.Rating)
		assert.NotNil(t, err)
	})
}

func TestIsCanGiveRating(t *testing.T) {
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
	transactionRepo = transaction.NewTransactionRepository(db)
	ratingRepo = rating.NewRatingRepository(db)

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
		Status:        "active",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	//CREATE TRANSACTION
	dummyTransaction := models.Transaction{
		PartnerID:  1,
		UserID:     1,
		Buffet:     false,
		Quantity:   1,
		Latitude:   100,
		Longtitude: 100,
		Distance:   1,
		Status:     "PAID",
	}
	db.Create(&dummyTransaction)

	dummyRating := models.Rating{
		TransactionID: 1,
		PartnerID:     1,
		UserID:        1,
		Rating:        5,
	}
	db.Create(&dummyRating)

	t.Run("IsCanGiveRating", func(t *testing.T) {
		dummyTransaction := models.Transaction{
			PartnerID:  1,
			UserID:     1,
			Buffet:     false,
			Quantity:   1,
			Latitude:   100,
			Longtitude: 100,
			Distance:   1,
			Status:     "CONFIRM",
		}
		db.Create(&dummyTransaction)

		res, _ := ratingRepo.IsCanGiveRating(int(dummyTransaction.UserID), int(dummyTransaction.ID))
		assert.Equal(t, true, res)
	})

	t.Run("IsCanGiveRating failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Rating{})
		var mockRating models.Rating
		mockRating.TransactionID = 99
		mockRating.PartnerID = 99
		mockRating.UserID = 1
		mockRating.Rating = 5

		_, err := ratingRepo.IsCanGiveRating(1, 1)
		// assert.Equal(t, 5, res.Rating)
		assert.NotNil(t, err)
	})
}
