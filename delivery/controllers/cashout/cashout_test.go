package cashout_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/cashout"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/xendit/xendit-go"
	"golang.org/x/crypto/bcrypt"
)

var JwtToken string

func TestCashout(t *testing.T) {
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

	t.Run("cashout bad request 4", func(t *testing.T) {
		// err := godotenv.Load()

		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }
		// xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

		e := echo.New()
		e.Validator = &cashout.CashoutValidator{Validator: validator.New()}
		bodyReq, _ := json.Marshal(cashout.CashoutRequest{

			BankCode:          "MANDIRI",
			AccountHolderName: "test",
			AccountNumber:     "1",
			Amount:            1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cashouts")

		cashoutController := cashout.NewCashoutController(mockCashout{})
		middleware.JWT([]byte(constants.JWT_SECRET_KEY))(cashoutController.Cashout)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("cashout success", func(t *testing.T) {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
		xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

		e := echo.New()
		e.Validator = &cashout.CashoutValidator{Validator: validator.New()}
		bodyReq, _ := json.Marshal(cashout.CashoutRequest{

			BankCode:          "MANDIRI",
			AccountHolderName: "test",
			AccountNumber:     "1",
			Amount:            1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cashouts")

		cashoutController := cashout.NewCashoutController(mockCashout{})
		middleware.JWT([]byte(constants.JWT_SECRET_KEY))(cashoutController.Cashout)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("cashout bad request 1", func(t *testing.T) {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
		xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

		e := echo.New()
		e.Validator = &cashout.CashoutValidator{Validator: validator.New()}
		bodyReq, _ := json.Marshal(cashout.CashoutRequest{

			// BankCode:          "MANDIRI",
			AccountHolderName: "test",
			AccountNumber:     "1",
			Amount:            1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cashouts")

		cashoutController := cashout.NewCashoutController(mockCashout{})
		middleware.JWT([]byte(constants.JWT_SECRET_KEY))(cashoutController.Cashout)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("cashout bad request 2", func(t *testing.T) {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
		xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

		e := echo.New()
		e.Validator = &cashout.CashoutValidator{Validator: validator.New()}
		bodyReq, _ := json.Marshal(cashout.CashoutRequest{

			BankCode:          "MANDIRI",
			AccountHolderName: "test",
			AccountNumber:     "1",
			Amount:            1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cashouts")

		cashoutController := cashout.NewCashoutController(mockFalseCashout{})
		middleware.JWT([]byte(constants.JWT_SECRET_KEY))(cashoutController.Cashout)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("cashout bad request 3", func(t *testing.T) {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
		xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

		e := echo.New()
		e.Validator = &cashout.CashoutValidator{Validator: validator.New()}
		bodyReq, _ := json.Marshal(cashout.CashoutRequest{

			BankCode:          "MANDIRI",
			AccountHolderName: "test",
			AccountNumber:     "1",
			Amount:            1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cashouts")

		cashoutController := cashout.NewCashoutController(mockCashout2{})
		middleware.JWT([]byte(constants.JWT_SECRET_KEY))(cashoutController.Cashout)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("cashout bad request 5", func(t *testing.T) {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
		xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

		e := echo.New()
		e.Validator = &cashout.CashoutValidator{Validator: validator.New()}
		bodyReq, _ := json.Marshal(cashout.CashoutRequest{

			BankCode:          "MANDIRI",
			AccountHolderName: "test",
			AccountNumber:     "1",
			Amount:            1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cashouts")

		cashoutController := cashout.NewCashoutController(mockFalseCashout2{})
		middleware.JWT([]byte(constants.JWT_SECRET_KEY))(cashoutController.Cashout)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})
}

//==========================
//MOCK CASHOUT
//==========================
type mockCashout struct{}

func (m mockCashout) Cashout(cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockCashout) History(userID int) ([]models.Cashout, error) {
	return []models.Cashout{
		{
			UserID: 1,
		},
	}, nil
}

func (m mockCashout) CheckBalance(userID int) (models.User, error) {
	return models.User{
		Email:   "test@gmail.com",
		Balance: 2000,
	}, nil
}

func (m mockCashout) CallbackSuccess(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockCashout) CallbackFailed(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

//==========================
//MOCK CASHOUT2
//==========================
type mockCashout2 struct{}

func (m mockCashout2) Cashout(cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockCashout2) History(userID int) ([]models.Cashout, error) {
	return []models.Cashout{
		{
			UserID: 1,
		},
	}, nil
}

func (m mockCashout2) CheckBalance(userID int) (models.User, error) {
	return models.User{
		Email:   "test@gmail.com",
		Balance: 500,
	}, nil
}

func (m mockCashout2) CallbackSuccess(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockCashout2) CallbackFailed(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

//==========================
//MOCK FALSE CASHOUT
//==========================
type mockFalseCashout struct{}

func (m mockFalseCashout) Cashout(cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockFalseCashout) History(userID int) ([]models.Cashout, error) {
	return []models.Cashout{
		{
			UserID: 1,
		},
	}, nil
}

func (m mockFalseCashout) CheckBalance(userID int) (models.User, error) {
	return models.User{
		Email:   "test@gmail.com",
		Balance: 2000,
	}, errors.New("FAILED")
}

func (m mockFalseCashout) CallbackSuccess(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockFalseCashout) CallbackFailed(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

//==========================
//MOCK FALSE CASHOUT2
//==========================
type mockFalseCashout2 struct{}

func (m mockFalseCashout2) Cashout(cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, errors.New("FAILED")
}

func (m mockFalseCashout2) History(userID int) ([]models.Cashout, error) {
	return []models.Cashout{
		{
			UserID: 1,
		},
	}, nil
}

func (m mockFalseCashout2) CheckBalance(userID int) (models.User, error) {
	return models.User{
		Email:   "test@gmail.com",
		Balance: 2000,
	}, nil
}

func (m mockFalseCashout2) CallbackSuccess(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
	}, nil
}

func (m mockFalseCashout2) CallbackFailed(extID string, cashout models.Cashout) (models.Cashout, error) {
	return models.Cashout{
		UserID: 1,
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
