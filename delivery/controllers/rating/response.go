package rating

type RatingResponse struct {
	TransactionID int `json:"transaction_id"`
	PartnerID  int    `json:"partner_id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}