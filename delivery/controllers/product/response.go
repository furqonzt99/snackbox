package product

import "github.com/furqonzt99/snackbox/models"

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

type PartnerResponse struct {
	Title       string  `json:"title"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
