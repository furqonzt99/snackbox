package partner

import (
	"net/http"

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
		apply.BussinessName = partnerReq.Bussiness_Name
		apply.Description = partnerReq.Description
		apply.Latitude = partnerReq.Latitude
		apply.Longtitude = partnerReq.Longtitude
		apply.LegalDocument = partnerReq.Legal_Document

		res, err := p.Repo.RequestPartner(apply)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(res))
	}
}
