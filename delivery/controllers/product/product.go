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

		productId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		updateProduct, err := p.Repo.FindProduct(productId, userJwt.UserID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		var product UpdateProductRequestFormat
		c.Bind(&product)

		if err2 := c.Validate(product); err2 != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		updateProduct.Title = product.Title
		updateProduct.Type = product.Type
		updateProduct.Description = product.Description
		updateProduct.Price = product.Price

		_, err3 := p.Repo.AddProduct(updateProduct)
		if err3 != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p ProductController) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		userJwt, _ := middlewares.ExtractTokenUser(c)

		productId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		err = p.Repo.DeleteProduct(productId, userJwt.UserID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())

	}
}

func (p ProductController) GetAllProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())

		}

		perpage, err := strconv.Atoi(c.QueryParam("perpage"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		search := c.QueryParam("search")

		if page == 0 {
			page = 1
		}

		if perpage == 0 {
			perpage = 10
		}

		offset := (page - 1) * perpage

		allProduct, _ := p.Repo.GetAllProduct(offset, perpage, search)

		productData := []GetProductWithPartnerResponse{}
		for _, item := range allProduct {
			productData = append(productData, GetProductWithPartnerResponse{
				Id:          item.ID,
				PartnerID:   item.PartnerID,
				Title:       item.Title,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		return c.JSON(http.StatusOK, common.PaginationResponse(page, perpage, productData))
	}
}
