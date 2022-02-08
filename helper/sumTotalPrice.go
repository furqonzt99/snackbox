package helper

import "github.com/xendit/xendit-go"

func SumTotalPrice(items []xendit.InvoiceItem) (totalPrice float64) {
	for _, item := range items {
		totalPrice += (item.Price * float64(item.Quantity))
	}

	return totalPrice
}