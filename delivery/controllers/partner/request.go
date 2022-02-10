package partner

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type PartnerUserRequestFormat struct {
	BussinessName string  `json:"bussiness_name" form:"name" validate:"required"`
	Description    string  `json:"description" form:"description" validate:"required"`
	Latitude       float64 `json:"latitude" form:"latitude" validate:"required"`
	Longtitude     float64 `json:"longtitude" form:"longtitude" validate:"required"`
	Address string  `json:"address" form:"address" validate:"required"`
	City string  `json:"city" form:"city" validate:"required"`
}

type UploadDocumentRequest struct {
	LegalDocument string  `json:"legal_document" form:"legal_document" validate:"required"`
}

type PartnerValidator struct {
	Validator *validator.Validate
}

func (cv *PartnerValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
