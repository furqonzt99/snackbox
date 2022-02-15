package transaction_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/transaction"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var JwtToken string //TOKEN FROM LOGIN
func TestTransaction(t *testing.T) {
	t.Run("Test Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("success transaction", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/order")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Order)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("transaction bad request validate", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			// Quantity:   1,
			Date:       "2022-02-30",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/order")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Order)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("transaction not found", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-30",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/order")

		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Order)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})

	t.Run(`transaction bad request "reservate 3 days before the event time"`, func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-16",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/order")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Order)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "you must reservate 3 days before the event time!", responses.Message)
	})

	t.Run("transaction bad request order", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-19",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/order")

		transactionController := transaction.NewTransactionController(mockFalseTransaction2{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Order)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})
}

func TestTransactionCallback(t *testing.T) {
	t.Run("callback success", func(t *testing.T) {
		e := echo.New()

		bodyReq, _ := json.Marshal(common.TransactionCallbackRequest{
			ExternalID: "1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/transactions/callback")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		transactionController.Callback(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("callback not found", func(t *testing.T) {
		e := echo.New()

		bodyReq, _ := json.Marshal(common.TransactionCallbackRequest{
			ExternalID: "1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/transactions/callback")

		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		transactionController.Callback(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
}

func TestAcceptTransaction(t *testing.T) {
	t.Run("Test Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("accept transaction success", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Accept)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("accept transaction badrequest param", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("a")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Accept)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("accept transaction not found acecpt repo", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Accept)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
}

func TestRejectTransaction(t *testing.T) {
	t.Run("Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("reject transaction success", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Reject)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("reject transaction badrequest param", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("a")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Reject)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("reject transaction not found reject repo", func(t *testing.T) {

		e := echo.New()
		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-22",
			Time:       "09:00:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Reject)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})

}

func TestSendTransaction(t *testing.T) {
	t.Run("Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("send transaction success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Send)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("send transaction badrequest param", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("a")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Send)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("send transaction err Repo.Send", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/send")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Send)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
}

func TestConfirmTransaction(t *testing.T) {
	t.Run("Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("confirm transaction success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/confirm")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Confirm)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("confirm transaction badrequest param", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/confirm")
		context.SetParamNames("id")
		context.SetParamValues("a")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Confirm)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("confirm transaction err Repo.Confirm", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/confirm")
		context.SetParamNames("id")
		context.SetParamValues("1")
		//yoga
		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Confirm)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
}

func TestGetAllTransaction(t *testing.T) {
	t.Run("Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("get all transaction success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.GetAll)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})
}

func TestGetOneTransaction(t *testing.T) {
	t.Run("Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("getone transaction success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.GetOne)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("confirm transaction badrequest param", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id")
		context.SetParamNames("id")
		context.SetParamValues("a")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.GetOne)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("confirm transaction err Repo.GetOneForUser", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")
		//yoga
		transactionController := transaction.NewTransactionController(mockFalseTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.GetOne)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
}

func TestShippingTransaction(t *testing.T) {
	t.Run("Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	t.Run("shipping success", func(t *testing.T) {

		e := echo.New()

		e.Validator = &transaction.TransactionValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(transaction.ShippingCostRequest{
			PartnerID:  1,
			Latitude:   100,
			Longtitude: 100,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))

		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := transaction.NewTransactionController(mockTransaction{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.Shipping)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	// t.Run("confirm transaction badrequest param", func(t *testing.T) {

	// 	e := echo.New()

	// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
	// 	res := httptest.NewRecorder()

	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

	// 	context := e.NewContext(req, res)
	// 	context.SetPath("/transactions/:id")
	// 	context.SetParamNames("id")
	// 	context.SetParamValues("a")

	// 	transactionController := transaction.NewTransactionController(mockTransaction{})
	// 	if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.GetOne)(context); err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}
	// 	var responses common.ResponseSuccess

	// 	json.Unmarshal([]byte(res.Body.Bytes()), &responses)
	// 	assert.Equal(t, "Bad Request", responses.Message)
	// })

	// t.Run("confirm transaction err Repo.GetOneForUser", func(t *testing.T) {

	// 	e := echo.New()

	// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
	// 	res := httptest.NewRecorder()

	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

	// 	context := e.NewContext(req, res)
	// 	context.SetPath("/transactions/:id/reject")
	// 	context.SetParamNames("id")
	// 	context.SetParamValues("1")
	// 	//yoga
	// 	transactionController := transaction.NewTransactionController(mockFalseTransaction{})
	// 	if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(transactionController.GetOne)(context); err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}
	// 	var responses common.ResponseSuccess

	// 	json.Unmarshal([]byte(res.Body.Bytes()), &responses)
	// 	assert.Equal(t, "Not Found", responses.Message)
	// })
}

//======================
//MOCK TRANSACTION REPOSITORY
//======================
type mockTransaction struct{}

func (m mockTransaction) Order(transaction models.Transaction, email string, products []int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
		Products: []models.Product{
			{
				Title: "bakso",
			},
		},
	}, nil
}

func (m mockTransaction) Accept(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockTransaction) Reject(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockTransaction) Send(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockTransaction) Confirm(trxID, userID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockTransaction) GetAllForPartner(partnerID int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:    1,
			PartnerID: 2,
		},
	}, nil
}

func (m mockTransaction) GetAllForUser(userID int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:    1,
			PartnerID: 2,
			Products: []models.Product{
				{
					Title: "bakso",
				},
			},
		},
	}, nil
}

func (m mockTransaction) GetOneForUser(trxID, userID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
		Products: []models.Product{
			{
				Title: "bakso",
			},
		},
	}, nil
}

func (m mockTransaction) GetOneForPartner(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
		Products: []models.Product{
			{
				Title: "bakso",
			},
		},
	}, nil
}

func (m mockTransaction) GetPartnerFromProduct(productID int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "test",
	}, nil
}

func (m mockTransaction) Callback(invId string, transaction models.Transaction, refund float64) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockTransaction) GetDistance(partnerID int, latitude, longtitude float64) (float64, error) {
	return 1, nil
}

//======================
//MOCK FALSE TRANSACTION REPOSITORY
//======================
type mockFalseTransaction struct{}

func (m mockFalseTransaction) Order(transaction models.Transaction, email string, products []int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction) Accept(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, errors.New("FAILED")
}

func (m mockFalseTransaction) Reject(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, errors.New("FAILED")
}

func (m mockFalseTransaction) Send(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, errors.New("FAILED")
}

func (m mockFalseTransaction) Confirm(trxID, userID int) (models.Transaction, error) {
	return models.Transaction{}, errors.New("FAILED")
}

func (m mockFalseTransaction) GetAllForPartner(partnerID int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:    1,
			PartnerID: 2,
		},
	}, nil
}

func (m mockFalseTransaction) GetAllForUser(userID int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:    1,
			PartnerID: 2,
		},
	}, errors.New("FAILED")
}

func (m mockFalseTransaction) GetOneForUser(trxID, userID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, errors.New("FAILED")
}

func (m mockFalseTransaction) GetOneForPartner(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction) GetPartnerFromProduct(productID int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "test",
	}, errors.New("FAILED")
}

func (m mockFalseTransaction) Callback(invId string, transaction models.Transaction, refund float64) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, errors.New("FAILED")

}
func (m mockFalseTransaction) GetDistance(partnerID int, latitude, longtitude float64) (float64, error) {
	return 1, errors.New("FAILED")
}

//======================
//MOCK FALSE TRANSACTION REPOSITORY2
//======================
type mockFalseTransaction2 struct{}

func (m mockFalseTransaction2) Order(transaction models.Transaction, email string, products []int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, errors.New("FAILED")
}

func (m mockFalseTransaction2) Accept(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) Reject(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) Send(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) Confirm(trxID, userID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) GetAllForPartner(partnerID int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:    1,
			PartnerID: 2,
		},
	}, nil
}

func (m mockFalseTransaction2) GetAllForUser(userID int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:    1,
			PartnerID: 2,
		},
	}, nil
}

func (m mockFalseTransaction2) GetOneForUser(trxID, userID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) GetOneForPartner(trxID, partnerID int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) GetPartnerFromProduct(productID int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "test",
	}, nil
}

func (m mockFalseTransaction2) Callback(invId string, transaction models.Transaction, refund float64) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
	}, nil
}

func (m mockFalseTransaction2) GetDistance(partnerID int, latitude, longtitude float64) (float64, error) {
	return 1, errors.New("FAILED")
}

//======================
//MOCK USER REPOSITORY
//======================
type mockUserRepository struct{}

func (m mockUserRepository) Register(newUser models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	return models.User{
		Email:    newUser.Email,
		Password: string(hash),
		Name:     newUser.Name,
		Address:  newUser.Address,
		City:     newUser.City,
	}, nil
}

func (m mockUserRepository) Login(email string) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return models.User{
		Email:    "test@gmail.com",
		Password: string(hash),
	}, nil
}

func (m mockUserRepository) Get(userid int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return models.User{
		Email:    "test@gmail.com",
		Password: string(hash),
		Name:     "tester"}, nil
}

func (m mockUserRepository) Update(newUser models.User, userId int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{
		Email:    "test2@gmail.com",
		Password: string(hash),
		Name:     "tester2",
	}, nil
}

func (m mockUserRepository) Delete(userId int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{
		Email:    "test2@gmail.com",
		Password: string(hash), Name: "tester2",
	}, nil
}
