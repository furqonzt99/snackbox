package partner

import (
	"time"

	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/controllers/rating"
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

type PartnerResponse struct {
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	LegalDocument string  `json:"legal_document"`
	Status        string  `json:"status"`
}

type GetPartnerResponse struct {
	ID            int     `json:"id"`
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	LegalDocument string  `json:"legal_document"`
	Status        string  `json:"status"`
}

type GetPartnerProductResponse struct {
	ID            int     `json:"id"`
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	Rating		  float64 `json:"rating"`
	Products      []product.ProductResponse
}

type GetPartnerRatingResponse struct {
	ID            int     `json:"id"`
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	Rating		  float64 `json:"rating"`
	Ratings      []rating.RatingResponse
}

type GetPartnerProfileResponse struct {
	ID            int     `json:"id"`
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	LegalDocument string  `json:"legal_document"`
	Status        string  `json:"status"`
	
}

type ReportResponse struct {
	CreateAt         time.Time              `json:"create_at"`
	InvoiceId        string                 `json:"invoice_id"`
	TotalTransaction float64                `json:"total_transaction"`
	Quantity         int                    `json:"quantity"`
	PaymentChannel   string                 `json:"payment_channel"`
	Status           string                 `json:"status"`
	Products         []ProductTitleResponse `json:"products"`
}

type ProductTitleResponse struct {
	Title string `json:"title"`
}

type PartnerData struct {
	ID            int     `json:"id"`
	BussinessName string  `json:"bussiness_name"`
	Description   string  `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longtitude"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	LegalDocument string  `json:"legal_document"`
	Status        string  `json:"status"`
	ApplyDate     string  `json:"apply_date"`
}
