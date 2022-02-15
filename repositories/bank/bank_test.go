package bank_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *user.UserRepository

func TestBank(t *testing.T) {
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

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Partner{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Rating{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.DetailTransaction{})
	db.AutoMigrate(&models.Cashout{})

	t.Run("Register User", func(t *testing.T) {
		var mockUser models.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test123"
		mockUser.Name = "tester"

		_, err := userRepo.Register(mockUser)
		assert.Nil(t, err)
		// assert.Equal(t, mockUser.Name, res.Name)
		// assert.Equal(t, 1, int(res.ID))
	})

	t.Run("Error Register User Duplicate Email", func(t *testing.T) {
		var mockUser models.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test123"
		mockUser.Name = "tester"

		_, err := userRepo.Register(mockUser)
		assert.NotNil(t, err)
	})
}
