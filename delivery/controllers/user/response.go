package user

import (
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/furqonzt99/snackbox/models"
)

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

type UserResponse struct {
	ID      uint    `json:"id"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	City    string  `json:"city"`
	Balance float64 `json:"balance"`
	Role    string  `json:"role"`
}

type UserProfileResponse struct {
	ID      uint    `json:"id"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
type UserProfileResponseWithPartner struct {
	ID      uint    `json:"id"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Partner partner.GetPartnerProfileResponse
}
