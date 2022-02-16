package user_test

import (
	"testing"

	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/seeder"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *user.UserRepository

func TestRegisterUser(t *testing.T) {
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

		res, _ := userRepo.Register(mockUser)
		// assert.Nil(t, err)
		assert.Equal(t, mockUser.Name, res.Name)
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

func TestLoginUser(t *testing.T) {
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

	dummyUser := models.User{
		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

	t.Run("Login User", func(t *testing.T) {
		var mockUser models.User
		mockUser.Email = "test@gmail.com"

		res, err := userRepo.Login(mockUser.Email)
		assert.Nil(t, err)
		assert.Equal(t, res.Email, mockUser.Email)
	})

	t.Run("Error Login User No Email", func(t *testing.T) {
		var mockUser models.User
		mockUser.Email = "test123@gmail.com"

		_, err := userRepo.Login(mockUser.Email)
		assert.NotNil(t, err)
	})
}

func TestGetUser(t *testing.T) {
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

	seeder.UserSeeder(db)

	t.Run("Get User", func(t *testing.T) {
		userId := 1
		res, _ := userRepo.Get(userId)
		// assert.Nil(t, err)
		assert.Equal(t, "User 2", res.Name)
	})

	t.Run("Error Get User No ID", func(t *testing.T) {
		userId := 100
		_, err := userRepo.Get(userId)
		assert.NotNil(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
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

	seeder.UserSeeder(db)

	t.Run("Update User ", func(t *testing.T) {
		var mockUser models.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test4321"
		mockUser.Name = "tester2"

		userId := 1

		res, _ := userRepo.Update(mockUser, userId)
		// assert.Nil(t, err)
		// assert.Equal(t, mockUser.Email, res.Email)
		// assert.Equal(t, mockUser.Password, res.Password)
		assert.Equal(t, "tester2", res.Name)
	})

	t.Run("Error Update User No ID", func(t *testing.T) {
		var mockUser models.User
		mockUser.Email = "test@gmail.com"

		userId := 100

		_, err := userRepo.Update(mockUser, userId)
		assert.NotNil(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
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

	seeder.UserSeeder(db)

	t.Run("Delete User", func(t *testing.T) {
		userId := 1
		res, _ := userRepo.Delete(userId)
		// assert.Nil(t, err)
		assert.Equal(t, "User 2", res.Name)
	})

	t.Run("Error Delete User No ID", func(t *testing.T) {
		userId := 100
		_, err := userRepo.Delete(userId)
		assert.NotNil(t, err)
	})
}
