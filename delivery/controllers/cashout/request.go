package cashout

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CashoutRequest struct {
	BankCode string `json:"bank_code" validate:"required"`
	AccountHolderName string `json:"account_holder_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}

type CashoutValidator struct {
	Validator *validator.Validate
}

func (cv *CashoutValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}