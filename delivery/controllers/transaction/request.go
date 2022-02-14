package transaction

import (
	"net/http"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TransactionRequest struct {
	Buffet bool `json:"buffet"`
	Quantity int `json:"quantity" validate:"required"`
	Date string `json:"date" validate:"required"`
	Time string `json:"time" validate:"required"`
	Latitude float64 `json:"latitude" validate:"required"`
	Longtitude float64 `json:"longtitude" validate:"required"`
	Products []int `json:"products" validate:"required"`
}

type ShippingCostRequest struct {
	PartnerID int `json:"partner_id" validate:"required"`
	Latitude float64 `json:"latitude" validate:"required"`
	Longtitude float64 `json:"longtitude" validate:"required"`
}

type TransactionValidator struct {
	Validator *validator.Validate
}

func (tv *TransactionValidator) Validate(i interface{}) error {
	if err := tv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return nil
}