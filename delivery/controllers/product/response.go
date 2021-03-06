package product

import (
	"github.com/furqonzt99/snackbox/models"
)

type GetUserResponseFormat struct {
	Message string      `json:"message"`
	Data    models.User `json:"data"`
}

type PutUserResponseFormat struct {
	Message string      `json:"message"`
	Data    models.User `json:"data"`
}

type DeleteUserResponseFormat struct {
	Message string `json:"message"`
}

type ProductResponse struct {
	Title       string  `json:"title"`
	Image		string	`json:"image"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type GetProductWithPartnerResponse struct {
	Id          uint    `json:"id"`
	PartnerID   uint    `json:"partner_id"`
	Title       string  `json:"title"`
	Image		string	`json:"image"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type GetPartnerResponse struct {
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
}
