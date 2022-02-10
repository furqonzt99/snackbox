package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPartnerPath(e *echo.Echo, partnerCtrl *partner.PartnerController) {

	e.POST("/partners/submission", partnerCtrl.ApplyPartner(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckUserRole)
	e.GET("/partners/submission", partnerCtrl.GetAllPartner(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckAdminRole)
	e.PUT("/partners/submission/:id/accept", partnerCtrl.AcceptPartner(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckAdminRole)
	e.PUT("/partners/submission/:id/reject", partnerCtrl.RejectPartner(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckAdminRole)
	e.GET("/partners/:id", partnerCtrl.GetPartner(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.POST("/partners/submission/upload", partnerCtrl.Upload, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckUserRole)
	e.GET("/partners/report", partnerCtrl.Report(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))

}
