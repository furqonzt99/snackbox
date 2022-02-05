package partner

import "github.com/furqonzt99/snackbox/models"

type RegisterUserResponseFormat struct {
	Message string        `json:"message"`
	Data    []models.User `json:"data"`
}

type LoginUserResponseFormat struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

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
	Bussiness_Name string  `json:"bussiness_Name"`
	Description    string  `json:"description"`
	Latitude       string  `json:"latitude"`
	Longtitude     string  `json:"longtitude"`
	Address        float64 `json:"address"`
	City           string  `json:"city"`
	Legal_Document string  `json:"legal_document"`
	Status         string  `json:"status"`
}
