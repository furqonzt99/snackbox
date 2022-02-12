package partner

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/labstack/echo/v4"
)

type PartnerController struct {
	Repo partner.PartnerInterface
}

func NewPartnerController(partner partner.PartnerInterface) *PartnerController {
	return &PartnerController{partner}
}

func (p PartnerController) ApplyPartner() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		partnerReq := PartnerUserRequestFormat{}
		c.Bind(&partnerReq)

		if err := c.Validate(partnerReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		var res models.Partner

		user, err := p.Repo.FindUserId(userJwt.UserID)

		var partnerDocument string
		if user.LegalDocument != "" {
			partnerDocument = fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, user.LegalDocument)
		}

		if err != nil {
			partnerData := models.Partner{
				UserID:        uint(userJwt.UserID),
				BussinessName: partnerReq.BussinessName,
				Description:   partnerReq.Description,
				Latitude:      partnerReq.Latitude,
				Longtitude:    partnerReq.Longtitude,
				Address:       partnerReq.Address,
				City:          partnerReq.City,
			}

			res, _ = p.Repo.ApplyPartner(partnerData)

			responseFormat := PartnerResponse{
				BussinessName: res.BussinessName,
				Description:   res.Description,
				Latitude:      res.Latitude,
				Longtitude:    res.Longtitude,
				Address:       res.Address,
				City:          res.City,
				LegalDocument: partnerDocument,
				Status:        res.Status,
			}

			return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))

		}

		if user.Status == "reject" {
			user.BussinessName = partnerReq.BussinessName
			user.Description = partnerReq.Description
			user.Latitude = partnerReq.Latitude
			user.Longtitude = partnerReq.Longtitude
			user.Address = partnerReq.Address
			user.City = partnerReq.City
			user.Status = "pending"

			res, _ = p.Repo.ApplyPartner(user)

			responseFormat := PartnerResponse{
				BussinessName: res.BussinessName,
				Description:   res.Description,
				Latitude:      res.Latitude,
				Longtitude:    res.Longtitude,
				Address:       res.Address,
				City:          res.City,
				LegalDocument: partnerDocument,
				Status:        res.Status,
			}
			return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))
		}

		responseFormat := PartnerResponse{
			BussinessName: user.BussinessName,
			Description:   user.Description,
			Latitude:      user.Latitude,
			Longtitude:    user.Longtitude,
			Address:       user.Address,
			City:          user.City,
			LegalDocument: partnerDocument,
			Status:        user.Status,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))
	}
}

func (p PartnerController) GetAllPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		res, err := p.Repo.GetAllPartner()
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		responseFormat := []GetPartnerResponse{}

		for _, data := range res {
			var partnerDocument string
			if data.LegalDocument != "" {
				partnerDocument = fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, data.LegalDocument)
			}
			responseFormat = append(responseFormat, GetPartnerResponse{
				ID:            int(data.ID),
				BussinessName: data.BussinessName,
				Description:   data.Description,
				Latitude:      data.Latitude,
				Longtitude:    data.Longtitude,
				Address:       data.Address,
				City:          data.City,
				LegalDocument: partnerDocument,
				Status:        data.Status,
			})
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))
	}
}

func (p PartnerController) AcceptPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		res, _ := p.Repo.FindPartnerId(partnerId)

		if res.Status == "reject" {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		err = p.Repo.AcceptPartner(res)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p PartnerController) RejectPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		res, _ := p.Repo.FindPartnerId(partnerId)

		if res.Status == "active" {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		err = p.Repo.RejectPartner(res)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p PartnerController) GetPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		partner, _ := p.Repo.GetPartner(partnerId)

		productItems := []product.ProductResponse{}
		for _, item := range partner.Products {
			productItems = append(productItems, product.ProductResponse{
				Title:       item.Title,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		response := GetPartnerProductResponse{
			ID:            int(partner.ID),
			BussinessName: partner.BussinessName,
			Description:   partner.Description,
			Latitude:      partner.Latitude,
			Longtitude:    partner.Longtitude,
			Address:       partner.Address,
			City:          partner.City,
			Products:      productItems,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}

}

func (pc PartnerController) Upload(c echo.Context) error {
	var requestUpload UploadDocumentRequest

	if err := c.Bind(&requestUpload); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	
	partner, err := pc.Repo.FindUserId(user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}
	
	if partner.Status == "active" || partner.Status == "pending" {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}
	
	file, err := c.FormFile("legal_document")
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

	const ACCEPTED_FILE_TYPE = "pdf"

	if kind.Extension != ACCEPTED_FILE_TYPE {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "extension must .pdf"))
	}

	fileID := strings.ReplaceAll(uuid.New().String(), "-", "")
	file.Filename = fmt.Sprint(fileID, ".", kind.Extension)

	if partner.LegalDocument != "" {
		if err := helper.GetObjectS3(partner.LegalDocument); err == nil {
			_ = helper.DeleteObjectS3(partner.LegalDocument)
		}
	}

	if err := helper.UploadObjectS3(file.Filename, src); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	const PENDING_STATUS = "pending"
	partnerData := models.Partner{
		LegalDocument: file.Filename,
		Status: PENDING_STATUS,
	}

	_, err = pc.Repo.UploadDocument(int(partner.ID), partnerData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}
