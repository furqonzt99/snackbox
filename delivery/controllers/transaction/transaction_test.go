package transaction_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/transaction"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

		bodyReq, _ := json.Marshal(transaction.TransactionRequest{
			Quantity:   1,
			Date:       "2022-02-11",
			Time:       "09:00",
			Latitude:   100,
			Longtitude: 100,
			Products:   []int{1},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/transactions/order")

	})
}

//======================
//MOCK TRANSACTION REPOSITORY
//======================

type mockTransaction struct{}

func (m mockTransaction) Order(transaction models.Transaction, email string, products []int) (models.Transaction, error) {
	return models.Transaction{
		UserID:    1,
		PartnerID: 2,
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
