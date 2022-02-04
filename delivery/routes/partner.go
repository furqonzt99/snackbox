package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPartnerPath(e *echo.Echo, partnerCtrl *partner.PartnerController) {

	e.POST("/partners", partnerCtrl.ApplyPartner(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))

}
