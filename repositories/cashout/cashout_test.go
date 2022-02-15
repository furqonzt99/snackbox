package cashout_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/cashout"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *user.UserRepository
var partnerRepo *partner.PartnerRepository
var cashoutRepo *cashout.CashoutRepository

func TestCashout(t *testing.T) {
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
	cashoutRepo = cashout.NewCashoutRepository(db)

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
		Balance:  1000,
	}
	userRepo.Register(dummyUser)

	//CREATE CASHOUT
	dummyCashout := models.Cashout{
		UserID: 1,
		Amount: 500,
	}
	db.Create(&dummyCashout)

	t.Run("cashout success", func(t *testing.T) {
		var mockCashout models.Cashout
		mockCashout.UserID = 1
		mockCashout.Amount = 300
		res, _ := cashoutRepo.Cashout(mockCashout)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount) //sudah
	})

	t.Run("cashout failed", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 10
		mockCashout2.Amount = 300

		_, err := cashoutRepo.Cashout(mockCashout2)
		assert.NotNil(t, err)
		// assert.Equal(t, float64(300), res.Amount) //sudah
	})

}

func TestHistory(t *testing.T) {
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
	cashoutRepo = cashout.NewCashoutRepository(db)

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
		Balance:  1000,
	}
	userRepo.Register(dummyUser)

	//CREATE CASHOUT
	dummyCashout := models.Cashout{
		UserID: 1,
		Amount: 500,
	}
	db.Create(&dummyCashout)

	t.Run("cashout success", func(t *testing.T) {

		res, _ := cashoutRepo.History(1)

		assert.Equal(t, float64(500), res[0].Amount)
	})

	t.Run("cashout failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Cashout{})
		_, err := cashoutRepo.History(1)

		assert.NotNil(t, err)
	})

}

func TestCallbackSuccess(t *testing.T) {
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
	cashoutRepo = cashout.NewCashoutRepository(db)

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
		Balance:  1000,
	}
	userRepo.Register(dummyUser)

	//CREATE CASHOUT
	dummyCashout := models.Cashout{
		UserID:     1,
		Amount:     500,
		ExternalID: "22",
	}
	db.Create(&dummyCashout)

	t.Run("CallbackSuccess success", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CallbackSuccess("22", mockCashout2)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})

	t.Run("CallbackSuccess failed 1", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CallbackSuccess("11", mockCashout2)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})

	t.Run("CallbackSuccess failed 2", func(t *testing.T) { //masih gagal
		var mockCashout2 models.Cashout
		mockCashout2.ID = 99
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CallbackSuccess("22", mockCashout2)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})
}

func TestCheckBalance(t *testing.T) {
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
	cashoutRepo = cashout.NewCashoutRepository(db)

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
		Balance:  1000,
	}
	userRepo.Register(dummyUser)

	//CREATE CASHOUT
	dummyCashout := models.Cashout{
		UserID:     1,
		Amount:     500,
		ExternalID: "22",
	}
	db.Create(&dummyCashout)

	t.Run("CheckBalance success", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CheckBalance(1)
		// assert.Nil(t, err)
		assert.Equal(t, float64(1000), res.Balance)
	})

	t.Run("CheckBalance success", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CheckBalance(3)
		// assert.Nil(t, err)
		assert.Equal(t, float64(0), res.Balance)
	})

}

func TestCallBackFailed(t *testing.T) {
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
	cashoutRepo = cashout.NewCashoutRepository(db)

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
		Balance:  1000,
	}
	userRepo.Register(dummyUser)

	//CREATE CASHOUT
	dummyCashout := models.Cashout{
		UserID:     1,
		Amount:     500,
		ExternalID: "22",
	}
	db.Create(&dummyCashout)

	t.Run("CheckBalance success", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CallbackFailed("22", mockCashout2)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})

	t.Run("CheckBalance failed 1", func(t *testing.T) {
		var mockCashout2 models.Cashout
		mockCashout2.UserID = 1
		mockCashout2.Amount = 300
		res, _ := cashoutRepo.CallbackFailed("99", mockCashout2)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})

	t.Run("CheckBalance failed 2", func(t *testing.T) {
		var mockCashout3 models.Cashout
		mockCashout3.UserID = 99
		mockCashout3.Amount = 300
		res, _ := cashoutRepo.CallbackFailed("22", mockCashout3)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})

	t.Run("CheckBalance failed 2", func(t *testing.T) {
		var mockCashout3 models.Cashout
		mockCashout3.UserID = 1
		mockCashout3.Amount = 300
		res, _ := cashoutRepo.CallbackFailed("22", mockCashout3)
		// assert.Nil(t, err)
		assert.Equal(t, float64(300), res.Amount)
	})

	// t.Run("CheckBalance success", func(t *testing.T) {
	// 	var mockCashout2 models.Cashout
	// 	mockCashout2.UserID = 1
	// 	mockCashout2.Amount = 300
	// 	res, _ := cashoutRepo.CheckBalance(3)
	// 	// assert.Nil(t, err)
	// 	assert.Equal(t, float64(0), res.Balance)
	// })

}
