package cashout

type CashoutResponse struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	IdempotenceKey string `json:"idempotence_key"`
	ExternalID string `json:"external_id"`
	BankCode string `json:"bank_code"`
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber string `json:"account_number"`
	Description string `json:"description"`
	Amount float64 `json:"amount"`
	Status string `json:"status"`
}