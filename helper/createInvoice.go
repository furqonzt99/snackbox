package helper

import (
	"os"

	"github.com/furqonzt99/snackbox/models"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(transaction models.Transaction, email string) (models.Transaction, error) {
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

	items := []xendit.InvoiceItem{}

	for _, product := range transaction.Products {
		items = append(items, xendit.InvoiceItem{
			Name:     product.Title,
			Price:    product.Price,
			Quantity: transaction.Quantity,
			Category: product.Type,
		})
	}

	transaction.TotalPrice = SumTotalPrice(items)

	data := invoice.CreateParams{
		ExternalID:      transaction.InvoiceID,
		Amount:          transaction.TotalPrice,
		Description:     "SnackBox Invoice " + transaction.InvoiceID + " for " + email,
		PayerEmail:      email,
		Items:           items,
	}

	resp, err := invoice.Create(&data)
	if err != nil {
		return transaction, err
	}

	transactionSuccess := models.Transaction{
		PaymentUrl:     resp.InvoiceURL,
		TotalPrice: resp.Amount,
	}

	return transactionSuccess, nil
}