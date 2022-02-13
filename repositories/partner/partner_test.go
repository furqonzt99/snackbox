package partner_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"

	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	usr "github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *usr.UserRepository
var partnerRepo *partner.PartnerRepository

func TestApplyPartner(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	t.Run("apply as partner", func(t *testing.T) {
		var mockPartner models.Partner
		mockPartner.UserID = 1
		mockPartner.BussinessName = "test"

		res, err := partnerRepo.ApplyPartner(mockPartner)
		assert.Nil(t, err)
		assert.Equal(t, 1, int(res.ID))
	})

	t.Run("apply as partner err", func(t *testing.T) {
		var mockPartner models.Partner
		mockPartner.UserID = 1
		mockPartner.BussinessName = "test"

		_, err := partnerRepo.ApplyPartner(mockPartner)
		assert.NotNil(t, err)
		// assert.Equal(t, 1, int(res.ID))
	})
}

func TestGetAllPartner(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	t.Run("get all partner", func(t *testing.T) {
		//APPLY PARTNER
		dummyPartner := models.Partner{
			UserID:        1,
			BussinessName: "partner1",
		}
		partnerRepo.ApplyPartner(dummyPartner)

		res, err := partnerRepo.GetAllPartner()
		assert.Nil(t, err)
		assert.Equal(t, "partner1", res[0].BussinessName)
	})

	t.Run("apply as partner failed", func(t *testing.T) {
		db.Migrator().DropTable(&models.Partner{})
		res, _ := partnerRepo.GetAllPartner()
		// assert.NotNil(t, err)
		assert.Equal(t, []models.Partner([]models.Partner(nil)), res)
	})

}

func TestFindPartner(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	t.Run("Find Partner Id", func(t *testing.T) {
		res, err := partnerRepo.FindPartnerId(1)
		assert.Nil(t, err)
		assert.Equal(t, "partner1", res.BussinessName)
	})

	t.Run("apply as partner err", func(t *testing.T) {
		db.Migrator().DropTable(&models.Partner{})
		res, _ := partnerRepo.FindPartnerId(1)
		// assert.NotNil(t, err)
		assert.Equal(t, "", res.BussinessName)
	})

}

func TestFindUserId(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	t.Run("Find User Id", func(t *testing.T) {
		res, err := partnerRepo.FindUserId(1)
		assert.Nil(t, err)
		assert.Equal(t, "partner1", res.BussinessName)
	})

	t.Run("apply as partner err", func(t *testing.T) {
		db.Migrator().DropTable(&models.Partner{})
		res, _ := partnerRepo.FindUserId(1)
		assert.Equal(t, "", res.BussinessName)
	})

}

func TestAcceptPartner(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	dummyPartner := models.Partner{
		// UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	t.Run("accept partner", func(t *testing.T) {
		partner := models.Partner{}
		partner.UserID = 1
		partner.BussinessName = "partner1"
		partner.Status = "pending"
		err := partnerRepo.AcceptPartner(partner)
		assert.Nil(t, err)

	})

	t.Run("accept partner", func(t *testing.T) {
		partner := models.Partner{}
		partner.UserID = 1
		partner.BussinessName = "partner1"
		partner.Status = "pending"
		db.Migrator().DropTable(&models.Partner{})
		err := partnerRepo.AcceptPartner(partner)
		assert.NotNil(t, err)

	})

}

func TestRejectPartner(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	dummyPartner := models.Partner{
		// UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	t.Run("reject partner", func(t *testing.T) {
		partner := models.Partner{}
		partner.UserID = 1
		partner.BussinessName = "partner1"
		partner.Status = "pending"
		err := partnerRepo.RejectPartner(partner)
		assert.Nil(t, err)

	})

	t.Run("reject partner failed", func(t *testing.T) {
		partner := models.Partner{}
		partner.UserID = 1
		partner.BussinessName = "partner1"
		partner.Status = "pending"
		db.Migrator().DropTable(&models.Partner{})
		err := partnerRepo.RejectPartner(partner)
		assert.NotNil(t, err)

	})

}

func TestGetAllPartnerProduct(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	t.Run("get all partner product", func(t *testing.T) {
		res, _ := partnerRepo.GetAllPartnerProduct()
		assert.Equal(t, "partner1", res[0].BussinessName)

	})

	t.Run("get all partner product", func(t *testing.T) {
		db.Migrator().DropTable(&models.Partner{})
		res, _ := partnerRepo.GetAllPartnerProduct()
		assert.Equal(t, []models.Partner([]models.Partner(nil)), res)

	})

}

func TestGetPartner(t *testing.T) {
	configTest = config.GetConfig()
	db = utils.InitDB(configTest)

	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Partner{})
	db.Migrator().DropTable(&models.Product{})
	db.Migrator().DropTable(&models.Rating{})
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.DetailTransaction{})
	db.Migrator().DropTable(&models.Cashout{})

	userRepo = usr.NewUserRepo(db)
	partnerRepo = partner.NewPartnerRepo(db)

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

	dummyPartner := models.Partner{
		UserID:        1,
		BussinessName: "partner1",
		Status:        "pending",
	}
	partnerRepo.ApplyPartner(dummyPartner)

	t.Run("get partner", func(t *testing.T) {
		res, _ := partnerRepo.GetPartner(1)
		assert.Equal(t, "partner1", res.BussinessName)

	})

	t.Run("get partner", func(t *testing.T) {
		db.Migrator().DropTable(&models.Partner{})
		res, _ := partnerRepo.GetPartner(1)
		assert.Equal(t, "", res.BussinessName)

	})

}
