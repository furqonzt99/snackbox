package helper

import (
	"fmt"
	"time"

	"github.com/furqonzt99/snackbox/models"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(transaction models.Transaction, email string, balance float64) (models.Transaction, error) {

	items := []xendit.InvoiceItem{}

	for _, product := range transaction.Products {
		items = append(items, xendit.InvoiceItem{
			Name:     product.Title,
			Price:    product.Price,
			Quantity: transaction.Quantity,
			Category: product.Type,
		})
	}

	shippingCost := CalculateShippingCost(transaction.Distance)

	items = append(items, xendit.InvoiceItem{
		Name:     "Shipping Cost",
		Price:    shippingCost,
		Quantity: 1,
	})

	transaction.TotalPrice = SumTotalPrice(items) + shippingCost

	totalPay := transaction.TotalPrice - balance

	var transactionSuccess models.Transaction

	if totalPay <= 0 {
		transactionSuccess = models.Transaction{
			TotalPrice:     transaction.TotalPrice,
			PaymentChannel: "SboxPay",
			PaymentMethod:  "Sboxpay",
			PaidAt:         time.Now(),
			Status:         "PAID",
		}
	} else {
		data := invoice.CreateParams{
			ExternalID:      transaction.InvoiceID,
			Amount:          totalPay,
			Description:     "SnackBox Invoice " + transaction.InvoiceID + " for " + email + " split with SboxPay Rp" + fmt.Sprint(balance),
			PayerEmail:      email,
			Items:           items,
		}

		resp, err := invoice.Create(&data)
		if err != nil {
			return transaction, err
		}

		transactionSuccess = models.Transaction{
			PaymentUrl:     resp.InvoiceURL,
			TotalPrice: transaction.TotalPrice,
		}
	}

	return transactionSuccess, nil
}