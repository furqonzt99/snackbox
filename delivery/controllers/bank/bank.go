package bank

import (
	"net/http"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/repositories/bank"
	"github.com/labstack/echo/v4"
)

type BankController struct {
	Repo bank.BankInterface
}

func NewBankController(bank bank.BankInterface) *BankController {
	return &BankController{Repo: bank}
}

func (bc BankController) AvailableBanks(c echo.Context) error {
	data, err := bc.Repo.GetAvailableBanks()
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(data))
}