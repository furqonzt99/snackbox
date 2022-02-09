package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/furqonzt99/snackbox/models"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go/disbursement"
)

func PaymentCashout(data models.Cashout) (models.Cashout, error) {

	externalID := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
	IdempotenceKey := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
	description := fmt.Sprint("Cashout to ", data.AccountHolderName, " at ", time.Now())

	createData := disbursement.CreateParams{
		IdempotencyKey:    IdempotenceKey,
		ExternalID:        externalID,
		BankCode:          data.BankCode,
		AccountHolderName: data.AccountHolderName,
		AccountNumber:     data.AccountNumber,
		Description:       description,
		Amount:            data.Amount,
	}

	resp, err := disbursement.Create(&createData)
	if err != nil {
		return data, err
	}

	disbursementData := models.Cashout{
		UserID:            data.UserID,
		IdempotenceKey:    IdempotenceKey,
		ExternalID:        externalID,
		BankCode:          resp.BankCode,
		AccountHolderName: resp.AccountHolderName,
		AccountNumber:     createData.AccountNumber,
		Amount:            resp.Amount,
		Description:       description,
		Status:            resp.Status,
	}
	
	return disbursementData, nil
}