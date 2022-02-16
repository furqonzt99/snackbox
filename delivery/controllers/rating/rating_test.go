package rating_test

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
	"github.com/furqonzt99/snackbox/delivery/controllers/rating"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JwtToken string //TOKEN FROM LOGIN
func TestRatingCreate(t *testing.T) {
	//======================
	//LOGIN TEST
	//======================
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

	t.Run("test post rating success", func(t *testing.T) {
		e := echo.New()
		e.Validator = &rating.RatingValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(rating.PostRatingRequest{
			Rating: 5,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("1")

		ratingController := rating.NewRatingController(mockRating{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("test post rating error param", func(t *testing.T) {
		e := echo.New()
		e.Validator = &rating.RatingValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(rating.PostRatingRequest{
			Rating: 5,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("a")

		ratingController := rating.NewRatingController(mockRating{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("test post rating error validate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &rating.RatingValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(rating.PostRatingRequest{
			Rating: 0,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("1")

		ratingController := rating.NewRatingController(mockRating{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("test post rating success", func(t *testing.T) {
		e := echo.New()
		e.Validator = &rating.RatingValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(rating.PostRatingRequest{
			Rating: 5,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("1")

		ratingController := rating.NewRatingController(mockFalseRating{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})

	t.Run("test post rating success", func(t *testing.T) {
		e := echo.New()
		e.Validator = &rating.RatingValidator{Validator: validator.New()}

		bodyReq, _ := json.Marshal(rating.PostRatingRequest{
			Rating: 5,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyReq))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("1")

		ratingController := rating.NewRatingController(mockFalseRating3{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}
		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})
}

func TestGetByTrxID(t *testing.T) {
	t.Run("get transaction id success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("1")

		ratingController := rating.NewRatingController(mockRating{})
		ratingController.GetByTrxID(context)

		responses := common.ResponseSuccess{}

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("get transaction id bad request", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("a")

		ratingController := rating.NewRatingController(mockRating{})
		ratingController.GetByTrxID(context)

		responses := common.ResponseSuccess{}

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("get transaction id success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/ratings/:trxID")
		context.SetParamNames("trxID")
		context.SetParamValues("1")

		ratingController := rating.NewRatingController(mockFalseRating{})
		ratingController.GetByTrxID(context)

		responses := common.ResponseSuccess{}

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
}

//==========================
//MOCK RATING
//==========================
type mockRating struct{}

func (m mockRating) Create(rating models.Rating) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, nil

}

func (m mockRating) Update(models.Rating) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, nil
}

func (m mockRating) IsCanGiveRating(userId, transactionId int) (models.Transaction, error) {
	return models.Transaction{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:    1,
		PartnerID: 1,
		Buffet:    false,
	}, nil
}

func (m mockRating) GetByTrxID(trxID int) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, nil
}

//==========================
//MOCK FALSE RATING
//==========================
type mockFalseRating struct{}

func (m mockFalseRating) Create(rating models.Rating) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, nil

}

func (m mockFalseRating) Update(models.Rating) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, nil
}

func (m mockFalseRating) IsCanGiveRating(userId, transactionId int) (models.Transaction, error) {
	return models.Transaction{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:    1,
		PartnerID: 1,
		Buffet:    false,
	}, errors.New("FAILED")
}

func (m mockFalseRating) GetByTrxID(trxID int) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, errors.New("FAILED")
}

//==========================
//MOCK FALSE RATING3
//==========================
type mockFalseRating3 struct{}

func (m mockFalseRating3) Create(rating models.Rating) (models.Rating, error) {
	return models.Rating{}, errors.New("FAILED")

}

func (m mockFalseRating3) Update(models.Rating) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, errors.New("FAILED")
}

func (m mockFalseRating3) IsCanGiveRating(userId, transactionId int) (models.Transaction, error) {
	return models.Transaction{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:    1,
		PartnerID: 1,
		Buffet:    false,
	}, nil
}

func (m mockFalseRating3) GetByTrxID(trxID int) (models.Rating, error) {
	return models.Rating{
		Rating: 5,
	}, errors.New("FAILED")
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
