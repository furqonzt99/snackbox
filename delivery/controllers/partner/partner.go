package partner

import (
	"net/http"
	"strconv"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
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
		if err != nil {
			partnerData := models.Partner{
				UserID:        uint(userJwt.UserID),
				BussinessName: partnerReq.BussinessName,
				Description:   partnerReq.Description,
				Latitude:      partnerReq.Latitude,
				Longtitude:    partnerReq.Longtitude,
				Address:       partnerReq.Address,
				City:          partnerReq.City,
				LegalDocument: partnerReq.LegalDocument,
			}

			res, _ = p.Repo.ApplyPartner(partnerData)

			responseFormat := PartnerResponse{
				BussinessName: res.BussinessName,
				Description:   res.Description,
				Latitude:      res.Latitude,
				Longtitude:    res.Longtitude,
				Address:       res.Address,
				City:          res.City,
				LegalDocument: res.LegalDocument,
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
			user.LegalDocument = partnerReq.LegalDocument
			user.Status = "pending"

			res, _ = p.Repo.ApplyPartner(user)

			responseFormat := PartnerResponse{
				BussinessName: res.BussinessName,
				Description:   res.Description,
				Latitude:      res.Latitude,
				Longtitude:    res.Longtitude,
				Address:       res.Address,
				City:          res.City,
				LegalDocument: res.LegalDocument,
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
			LegalDocument: user.LegalDocument,
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
			responseFormat = append(responseFormat, GetPartnerResponse{
				ID:            int(data.ID),
				BussinessName: data.BussinessName,
				Description:   data.Description,
				Latitude:      data.Latitude,
				Longtitude:    data.Longtitude,
				Address:       data.Address,
				City:          data.City,
				LegalDocument: data.LegalDocument,
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
