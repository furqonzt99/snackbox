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

		var apply models.Partner
		apply.UserID = uint(userJwt.UserID)
		apply.BussinessName = partnerReq.BussinessName
		apply.Description = partnerReq.Description
		apply.Latitude = partnerReq.Latitude
		apply.Longtitude = partnerReq.Longtitude
		apply.LegalDocument = partnerReq.LegalDocument

		res, err := p.Repo.RequestPartner(apply)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(res))
	}
}

func (p PartnerController) GetAllPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		res, err := p.Repo.GetAllPartner()
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewNotFoundResponse())
		}

		if len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		return c.JSON(http.StatusOK, common.SuccessResponse(res))
	}
}

func (p PartnerController) DeletePartner() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		err := p.Repo.DeletePartner(userJwt.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
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
