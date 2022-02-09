package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/cashout"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterCashoutPath(e *echo.Echo, CashoutController *cashout.CashoutController) {

	e.POST("/cashouts", CashoutController.Cashout, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.GET("/cashouts", CashoutController.History, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.POST("/cashouts/callback", CashoutController.Callback, middlewares.CheckXHeaderToken)
}