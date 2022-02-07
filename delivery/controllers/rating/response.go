package rating

type RatingResponse struct {
	PartnerID  int    `json:"partner_id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}