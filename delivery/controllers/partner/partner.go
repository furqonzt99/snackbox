package partner

import (
	"net/http"
	"strconv"

	"github.com/furqonzt99/snackbox/delivery/common"
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

		var partner models.Partner
		partner.UserID = uint(userJwt.UserID)
		partner.BussinessName = partnerReq.BussinessName
		partner.Description = partnerReq.Description
		partner.Latitude = partnerReq.Latitude
		partner.Longtitude = partnerReq.Longtitude
		partner.Address = partnerReq.Address
		partner.City = partnerReq.City
		partner.LegalDocument = partnerReq.LegalDocument

		res, err := p.Repo.RequestPartner(partner)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "already exist"))
		}

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
}

func (p PartnerController) GetAllPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		res, _ := p.Repo.GetAllPartner()
		if len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
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

		partnerId, _ := strconv.Atoi(c.Param("id"))

		res, err := p.Repo.FindPartnerId(partnerId)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		err2 := p.Repo.AcceptPartner(res)
		if err2 != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p PartnerController) RejectPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, _ := strconv.Atoi(c.Param("id"))

		res, err := p.Repo.FindPartnerId(partnerId)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		err2 := p.Repo.RejectPartner(res)
		if err2 != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

// func (p PartnerController) GetAllPartnerProduct() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		partnerId, _ := strconv.Atoi(c.Param("id"))
// 		res, err := p.Repo.GetAllPartnerProduct()
// 		if err != nil {
// 			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
// 		}
// 		if len(res) == 0 {
// 			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
// 		}

// 		return c.JSON(http.StatusOK, common.SuccessResponse(res))
// 	}
// }
