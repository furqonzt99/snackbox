package cashout

import (
	"net/http"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/cashout"
	"github.com/labstack/echo/v4"
)

type CashoutController struct {
	Repo cashout.CashoutInterface
}

func NewCashoutController(cashout cashout.CashoutInterface) *CashoutController {
	return &CashoutController{Repo: cashout}
}

func (cc CashoutController) Cashout(c echo.Context) error {
	var requestCashout CashoutRequest

	c.Bind(&requestCashout)

	if err := c.Validate(&requestCashout); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())

	}

	user, _ := middlewares.ExtractTokenUser(c)

	userData, err := cc.Repo.CheckBalance(user.UserID)
	if err != nil {
		// return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())

	}

	if requestCashout.Amount > userData.Balance {
		// return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())

	}

	data := models.Cashout{
		UserID:            uint(user.UserID),
		BankCode:          requestCashout.BankCode,
		AccountHolderName: requestCashout.AccountHolderName,
		AccountNumber:     requestCashout.AccountNumber,
		Amount:            requestCashout.Amount,
	}

	cashout, err := helper.PaymentCashout(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		// return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "ERROR6"))

	}

	cashoutDB, err := cc.Repo.Cashout(cashout)
	if err != nil {
		// return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())

	}

	response := CashoutResponse{
		ID:                int(cashoutDB.ID),
		UserID:            int(cashoutDB.UserID),
		IdempotenceKey:    cashoutDB.IdempotenceKey,
		ExternalID:        cashoutDB.ExternalID,
		BankCode:          cashoutDB.BankCode,
		AccountHolderName: cashoutDB.AccountHolderName,
		AccountNumber:     cashoutDB.AccountNumber,
		Description:       cashoutDB.Description,
		Amount:            cashoutDB.Amount,
		Status:            cashoutDB.Status,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}

func (cc CashoutController) History(c echo.Context) error {
	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	cashouts, err := cc.Repo.History(user.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	response := []CashoutResponse{}

	for _, cashout := range cashouts {
		response = append(response, CashoutResponse{
			ID:                int(cashout.ID),
			UserID:            int(cashout.UserID),
			IdempotenceKey:    cashout.IdempotenceKey,
			ExternalID:        cashout.ExternalID,
			BankCode:          cashout.BankCode,
			AccountHolderName: cashout.AccountHolderName,
			AccountNumber:     cashout.AccountNumber,
			Description:       cashout.Description,
			Amount:            cashout.Amount,
			Status:            cashout.Status,
		})
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}

func (cc CashoutController) Callback(c echo.Context) error {

	var callbackRequest common.CashoutCallbackRequest
	if err := c.Bind(&callbackRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var data models.Cashout
	data.ExternalID = callbackRequest.ExternalID
	data.Amount = callbackRequest.Amount
	data.Status = callbackRequest.Status

	const STATUS_COMPLETED = "COMPLETED"

	var err error

	if callbackRequest.Status == STATUS_COMPLETED {
		_, err = cc.Repo.CallbackSuccess(callbackRequest.ExternalID, data)
	} else {
		_, err = cc.Repo.CallbackFailed(callbackRequest.ExternalID, data)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}
