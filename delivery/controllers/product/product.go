package product

import (
	"net/http"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/product"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	Repo product.ProductInterface
}

func NewProductController(product product.ProductInterface) *ProductController {
	return &ProductController{product}
}

func (p ProductController) AddProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		var product RegisterProductRequestFormat
		c.Bind(&product)

		if err := c.Validate(product); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		var add models.Product
		add.PartnerID = uint(userJwt.UserID)
		add.Title = product.Title
		add.Type = product.Type
		add.Description = product.Description
		add.Price = product.Price

		res, err := p.Repo.AddProduct(add)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		response := PartnerResponse{
			Title:       res.Title,
			Type:        res.Type,
			Description: res.Description,
			Price:       res.Price,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}
