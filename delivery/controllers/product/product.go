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

		var productReq RegisterProductRequestFormat
		c.Bind(&productReq)

		if err := c.Validate(productReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		var product models.Product
		product.PartnerID = uint(userJwt.PartnerID)
		product.Title = productReq.Title
		product.Type = productReq.Type
		product.Description = productReq.Description
		product.Price = productReq.Price

		res, err := p.Repo.AddProduct(product)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		response := ProductResponse{
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

		// partner := GetPartnerResponse{}
		// for _, item := range

		productWithPartner := []GetProductWithPartnerResponse{}
		for _, item := range allProduct {
			productWithPartner = append(productWithPartner, GetProductWithPartnerResponse{
				Id:          item.ID,
				Title:       item.Title,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
				// Partner: GetPartnerResponse{
				// 	BussinessName: item.Partner.BussinessName,
				// 	Description:   item.Partner.Description,
				// 	Latitude:      item.Partner.Latitude,
				// 	Longtitude:    item.Partner.Longtitude,
				// 	Address:       item.Partner.Address,
				// 	City:          item.Partner.City,
				// },
			})
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(productWithPartner))
	}
}

func (p ProductController) SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		search := c.QueryParam("search")

		product, err := p.Repo.SearchProduct(search)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewNotFoundResponse())
		}

		res := []ProductResponse{}
		for _, item := range product {
			res = append(res, ProductResponse{
				Title:       item.Title,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(res))
	}
}
