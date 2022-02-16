package bank_test

import (
	"log"
	"os"
	"testing"

	"github.com/furqonzt99/snackbox/repositories/bank"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/xendit/xendit-go"
	"gorm.io/gorm"
)

func TestBank(t *testing.T) {
	var db *gorm.DB
	bankRepo := bank.NewBankRepository(db)

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

	t.Run("GetAvailableBanks", func(t *testing.T) {

		res, _ := bankRepo.GetAvailableBanks()
		assert.Equal(t, "Bank Central Asia (BCA)", res[0].Name)
	})
}

func TestBankFailed(t *testing.T) {
	var bankRepositoryFail *bank.BankRepository
	err := godotenv.Load("./failedTest/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY_FAILED")

	t.Run("GetAvailableBanksFialed", func(t *testing.T) {

		_, err := bankRepositoryFail.GetAvailableBanks()
		assert.NotNil(t, err)
	})
}
