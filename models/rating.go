package models

type Rating struct {
	PartnerID uint
	UserID uint
	Rating int
	Comment string
	User User
	Partner Partner
} 