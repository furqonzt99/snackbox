package transaction

import "github.com/furqonzt99/snackbox/delivery/controllers/product"

type TransactionResponse struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	UserName string `json:"username"`
	PartnerID int `json:"partner_id"`
	InvoiceID string `json:"invoice_id"`
	Buffet bool `json:"buffet"`
	Quantity int `json:"quantity"`
	Latitude float64 `json:"latitude"`
	Longtitude float64 `json:"longtitude"`
	DateTime string `json:"datetime"`
	Distance float32 `json:"distance"`
	TotalPrice float64 `json:"total_price"`
	PaymentUrl string `json:"payment_url"`
	PaymentMethod string `json:"payment_method"`
	PaymentChannel string `json:"payment_channel"`
	PaidAt string `json:"paid_at"`
	Status string `json:"status"`
	Products []product.ProductResponse `json:"products"`
}