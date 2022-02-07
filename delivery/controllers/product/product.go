package product

import (
	"net/http"
	"strconv"

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
func (p ProductController) PutProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		productId, _ := strconv.Atoi(c.Param("id"))

		fond, err := p.Repo.FindProduct(productId, userJwt.UserID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		var product UpdateProductRequestFormat
		c.Bind(&product)

		if err2 := c.Validate(product); err2 != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		// var add models.Product
		// fond.PartnerID = uint(userJwt.UserID)
		fond.Title = product.Title
		fond.Type = product.Type
		fond.Description = product.Description
		fond.Price = product.Price

		_, err3 := p.Repo.AddProduct(fond)
		if err3 != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p ProductController) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		userJwt, _ := middlewares.ExtractTokenUser(c)

		productId, _ := strconv.Atoi(c.Param("id"))

		err := p.Repo.DeleteProduct(productId, userJwt.UserID)

		if err != nil {
			return c.JSON(http.StatusBadGateway, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())

	}
}

func (p ProductController) GetAllProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		allProduct, _ := p.Repo.GetAllProduct()

		if len(allProduct) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(allProduct))
	}
}
