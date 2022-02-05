package common

type JWTPayload struct {
	UserID int
	PartnerID int
	Email string
	Role string
}