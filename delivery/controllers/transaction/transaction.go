package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/transaction"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	Repo transaction.TransactionInterface
}

func NewTransactionController(repo transaction.TransactionInterface) *TransactionController {
	return &TransactionController{Repo: repo}
}

func (tc *TransactionController) Order(c echo.Context) error {
	var transactionRequest TransactionRequest

	c.Bind(&transactionRequest)

	if err := c.Validate(&transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := middlewares.ExtractTokenUser(c)

	//get partner id from product
	partner, err := tc.Repo.GetPartnerFromProduct(transactionRequest.Products[0])
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	// create invoiceID
	invoiceId := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))

	dateTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("%vT%vZ", transactionRequest.Date, transactionRequest.Time))
	threeDayLater := time.Now().AddDate(0, 0, 3)

	if threeDayLater.After(dateTime) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "you must reservate 3 days before the event time!"))
	}

	distance, err := tc.Repo.GetDistance(int(partner.ID), transactionRequest.Latitude, transactionRequest.Longtitude)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	transaction := models.Transaction{
		UserID:     uint(user.UserID),
		PartnerID:  uint(partner.ID),
		Buffet:     transactionRequest.Buffet,
		Quantity:   transactionRequest.Quantity,
		DateTime:   dateTime,
		Latitude:   transactionRequest.Latitude,
		Longtitude: transactionRequest.Longtitude,
		Distance:   distance,
		InvoiceID:  invoiceId,
	}

	transactionOrder, err := tc.Repo.Order(transaction, user.Email, transactionRequest.Products)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	productItems := []product.ProductResponse{}

	for _, data := range transactionOrder.Products {
		productItems = append(productItems, product.ProductResponse{
			Title:       data.Title,
			Type:        data.Type,
			Description: data.Description,
			Price:       data.Price,
		})
	}

	response := TransactionResponse{
		ID:             int(transactionOrder.ID),
		UserID:         int(transactionOrder.UserID),
		UserName:       transactionOrder.User.Name,
		PartnerID:      int(transactionOrder.PartnerID),
		InvoiceID:      transactionOrder.InvoiceID,
		Buffet:         transactionOrder.Buffet,
		Quantity:       transactionOrder.Quantity,
		Latitude:       transactionOrder.Latitude,
		Longtitude:     transaction.Longtitude,
		DateTime:       fmt.Sprint(transactionOrder.DateTime),
		Distance:       float32(transactionOrder.Distance),
		TotalPrice:     transactionOrder.TotalPrice,
		PaymentUrl:     transactionOrder.PaymentUrl,
		PaymentMethod:  transactionOrder.PaymentMethod,
		PaymentChannel: transactionOrder.PaymentChannel,
		PaidAt:         fmt.Sprint(transactionOrder.PaidAt),
		Status:         transactionOrder.Status,
		Products:       productItems,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}

func (tc TransactionController) Callback(c echo.Context) error {

	var callbackRequest common.TransactionCallbackRequest
	c.Bind(&callbackRequest)

	var data models.Transaction
	data.PaidAt, _ = time.Parse(time.RFC3339, callbackRequest.PaidAt)
	data.PaymentMethod = callbackRequest.PaymentMethod
	data.PaymentChannel = callbackRequest.PaymentChannel
	data.Status = callbackRequest.Status

	var refund float64

	if callbackRequest.Status == "PAID" {
		refund = 0
	} else {
		for _, item := range callbackRequest.Items {
			refund += item.Price
		}
	}

	_, err := tc.Repo.Callback(callbackRequest.ExternalID, data, refund)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) Accept(c echo.Context) error {

	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := middlewares.ExtractTokenUser(c)

	_, err = tc.Repo.Accept(trxID, user.PartnerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) Reject(c echo.Context) error {

	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := middlewares.ExtractTokenUser(c)

	_, err = tc.Repo.Reject(trxID, user.PartnerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) Send(c echo.Context) error {

	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := middlewares.ExtractTokenUser(c)

	_, err = tc.Repo.Send(trxID, user.PartnerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) Confirm(c echo.Context) error {

	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = tc.Repo.Confirm(trxID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) GetAll(c echo.Context) error {
	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var data []models.Transaction
	if user.PartnerID != 0 {
		data, _ = tc.Repo.GetAllForPartner(user.PartnerID)
	} else {
		data, _ = tc.Repo.GetAllForUser(user.UserID)
	}

	response := []TransactionResponse{}

	for _, trx := range data {

		productItems := []product.ProductResponse{}
		for _, item := range trx.Products {
			productItems = append(productItems, product.ProductResponse{
				Title:       item.Title,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		response = append(response, TransactionResponse{
			ID:             int(trx.ID),
			UserID:         int(trx.UserID),
			UserName:       trx.User.Name,
			PartnerID:      int(trx.PartnerID),
			InvoiceID:      trx.InvoiceID,
			Buffet:         trx.Buffet,
			Quantity:       trx.Quantity,
			Latitude:       trx.Latitude,
			Longtitude:     trx.Longtitude,
			DateTime:       fmt.Sprint(trx.DateTime),
			Distance:       float32(trx.Distance),
			TotalPrice:     trx.TotalPrice,
			PaymentUrl:     trx.PaymentUrl,
			PaymentMethod:  trx.PaymentMethod,
			PaymentChannel: trx.PaymentChannel,
			PaidAt:         fmt.Sprint(trx.PaidAt),
			Status:         trx.Status,
			Products:       productItems,
		})
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}

func (tc TransactionController) GetOne(c echo.Context) error {

	trxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	data := models.Transaction{}

	if user.PartnerID != 0 {
		data, err = tc.Repo.GetOneForPartner(trxID, user.PartnerID)
	} else {
		data, err = tc.Repo.GetOneForUser(trxID, user.UserID)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	response := []TransactionResponse{}

	productItems := []product.ProductResponse{}
	for _, item := range data.Products {
		productItems = append(productItems, product.ProductResponse{
			Title:       item.Title,
			Type:        item.Type,
			Description: item.Description,
			Price:       item.Price,
		})
	}

	response = append(response, TransactionResponse{
		ID:             int(data.ID),
		UserID:         int(data.UserID),
		UserName:       data.User.Name,
		PartnerID:      int(data.PartnerID),
		InvoiceID:      data.InvoiceID,
		Buffet:         data.Buffet,
		Quantity:       data.Quantity,
		Latitude:       data.Latitude,
		Longtitude:     data.Longtitude,
		DateTime:       fmt.Sprint(data.DateTime),
		Distance:       float32(data.Distance),
		TotalPrice:     data.TotalPrice,
		PaymentUrl:     data.PaymentUrl,
		PaymentMethod:  data.PaymentMethod,
		PaymentChannel: data.PaymentChannel,
		PaidAt:         fmt.Sprint(data.PaidAt),
		Status:         data.Status,
		Products:       productItems,
	})

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}

func (tc TransactionController) Shipping(c echo.Context) error {
	var shippingRequest ShippingCostRequest

	if err := c.Bind(&shippingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if err := c.Validate(&shippingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	distance, err := tc.Repo.GetDistance(shippingRequest.PartnerID, shippingRequest.Latitude, shippingRequest.Longtitude)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	cost := helper.CalculateShippingCost(distance)

	response := ShippingCostResponse{
		Distance: distance,
		Cost:     cost,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}
