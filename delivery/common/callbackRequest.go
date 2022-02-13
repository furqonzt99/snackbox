package common

type TransactionCallbackRequest struct {
	ExternalID string `json:"external_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentChannel string `json:"payment_channel"`
	PaidAt string `json:"paid_at"`
	Status string `json:"status"`
	Items []XenditItems `json:"items"`
}

type CashoutCallbackRequest struct {
	ExternalID string `json:"external_id"`
	Amount float64 `json:"amount"`
	Status string `json:"status"`
}

type XenditItems struct {
	Name string `json:"name"`
    Price float64 `json:"price"`
    Quantity int `json:"quantity"`
}