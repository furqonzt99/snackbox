package product

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/product"
	"github.com/google/uuid"
	"github.com/h2non/filetype"

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
			Image: 		 res.Image,
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

		page, _ := strconv.Atoi(c.QueryParam("page"))
		perpage, _ := strconv.Atoi(c.QueryParam("perpage"))
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
			var productImage string
			if item.Image != "" { 
				productImage = fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, item.Image)
			}
			productData = append(productData, GetProductWithPartnerResponse{
				Id:          item.ID,
				PartnerID:   item.PartnerID,
				Title:       item.Title,
				Image: 		 productImage,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		return c.JSON(http.StatusOK, common.PaginationResponse(page, perpage, productData))
	}
}

func (pc ProductController) Upload(c echo.Context) error {

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var requestUpload UpdateProductRequestFormat

	if err := c.Bind(&requestUpload); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	
	product, err := pc.Repo.FindProduct(productID, user.PartnerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}
	
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	defer src.Close()

	head := make([]byte, 261)
  	src.Read(head)

	kind, _ := filetype.Match(head)

	if !filetype.IsImage(head) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "file type must an image"))
	}

	fileID := strings.ReplaceAll(uuid.New().String(), "-", "")
	file.Filename = fmt.Sprint(fileID, ".", kind.Extension)

	if product.Image != "" {
		if err := helper.GetObjectS3(product.Image); err == nil {
			_ = helper.DeleteObjectS3(product.Image)
		}
	}

	if err := helper.UploadObjectS3(file.Filename, src); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	productData := models.Product{
		Image: file.Filename,
	}

	_, err = pc.Repo.UploadImage(productID, productData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}
