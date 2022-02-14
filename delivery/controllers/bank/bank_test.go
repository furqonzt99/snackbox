package bank_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/bank"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xendit/xendit-go"
)

func TestBank(t *testing.T) {
	t.Run("test GetAvailableBanks success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/banks")

		bankController := bank.NewBankController(mockBankRepository{})
		bankController.AvailableBanks(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("test GetAvailableBanks failed", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/banks")

		bankController := bank.NewBankController(mockFalseBankRepository{})
		bankController.AvailableBanks(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})
}

//MOCK BANK
type mockBankRepository struct{}

func (m mockBankRepository) GetAvailableBanks() ([]xendit.DisbursementBank, error) {
	return []xendit.DisbursementBank{{
		Name:            "Bank Mandiri",
		Code:            "MANDIRI",
		CanDisburse:     true,
		CanNameValidate: true,
	}}, nil
}

type mockFalseBankRepository struct{}

func (m mockFalseBankRepository) GetAvailableBanks() ([]xendit.DisbursementBank, error) {
	return []xendit.DisbursementBank{{
		Name:            "Bank Mandiri",
		Code:            "MANDIRI",
		CanDisburse:     true,
		CanNameValidate: true,
	}}, errors.New("failed")
}
