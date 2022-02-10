package product

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RegisterProductRequestFormat struct {
	Title       string  `json:"title" form:"title" validate:"required"`
	Type        string  `json:"type" form:"type" validate:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" validate:"required"`
}

type UpdateProductRequestFormat struct {
	Title       string  `json:"title" form:"title" validate:"required"`
	Type        string  `json:"type" form:"type" validate:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" validate:"required"`
}

type UploadProductRequestFormat struct {
	Image       string  `form:"title" validate:"required"`
}
type ProductValidator struct {
	Validator *validator.Validate
}

func (cv *ProductValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
