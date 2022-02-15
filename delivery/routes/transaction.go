package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/transaction"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterTransactionPath(e *echo.Echo, TransactionController *transaction.TransactionController) {

	e.POST("/transactions/order", TransactionController.Order, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckUserRole)
	e.POST("/transactions/callback", TransactionController.Callback, middlewares.CheckXHeaderToken)
	e.PUT("/transactions/:id/accept", TransactionController.Accept, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
	e.PUT("/transactions/:id/reject", TransactionController.Reject, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
	e.PUT("/transactions/:id/send", TransactionController.Send, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
	e.PUT("/transactions/:id/confirm", TransactionController.Confirm, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckUserRole)
	e.GET("/transactions", TransactionController.GetAll, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.GET("/transactions/:id", TransactionController.GetOne, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.GET("/transactions/shipping", TransactionController.Shipping, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
}