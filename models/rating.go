package models

type Rating struct {
	TransactionID uint
	PartnerID uint
	UserID uint
	Rating int
	Comment string
	Transaction Transaction
	User User
	Partner Partner
} 