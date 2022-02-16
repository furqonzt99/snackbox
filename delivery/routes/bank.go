package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/bank"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterBankPath(e *echo.Echo, BankController *bank.BankController)  {
	e.GET("/banks",BankController.AvailableBanks, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
}