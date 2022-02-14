package partner_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JwtToken string

func TestPartner(t *testing.T) {
	t.Run(
		"Test Login", func(t *testing.T) {
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
			// fmt.Println(response)
			JwtToken = response.Data.(string)
			// assert.Equal(t, "Successful Operation", response.Message)
			// assert.NotNil(t, JwtToken)
		})

	t.Run(
		"test Apply as Partner Success", func(t *testing.T) {
			e := echo.New()

			e.Validator = &partner.PartnerValidator{Validator: validator.New()}
			requestBody, _ := json.Marshal(partner.PartnerUserRequestFormat{
				BussinessName: "testPartner",
				Description:   "testPartner",
				Latitude:      100,
				Longtitude:    100,
				Address:       "testPartner",
				City:          "testPartner",
			})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

			context := e.NewContext(req, res)
			context.SetPath("/partners/submission")

			partnerController := partner.NewPartnerController(mockPartnerRepository{})
			if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.ApplyPartner())(context); err != nil {
				log.Fatal(err)
				return
			}

			var responses common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &responses)
			assert.Equal(t, "Successful Operation", responses.Message)
		})

	t.Run(
		"test Apply as Partner validation bad request", func(t *testing.T) {
			e := echo.New()

			e.Validator = &partner.PartnerValidator{Validator: validator.New()}
			requestBody, _ := json.Marshal(partner.PartnerUserRequestFormat{
				BussinessName: "",
				Description:   "testPartner",
				Latitude:      100,
				Longtitude:    100,
				Address:       "testPartner",
				City:          "testPartner",
			})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

			context := e.NewContext(req, res)
			context.SetPath("/partners/submission")

			partnerController := partner.NewPartnerController(mockPartnerRepository{})
			if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.ApplyPartner())(context); err != nil {
				log.Fatal(err)
				return
			}

			var responses common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &responses)
			assert.Equal(t, "Bad Request", responses.Message)
		})

	t.Run(
		"test Apply as Partner status reject", func(t *testing.T) {
			e := echo.New()

			e.Validator = &partner.PartnerValidator{Validator: validator.New()}
			requestBody, _ := json.Marshal(partner.PartnerUserRequestFormat{
				BussinessName: "testPartner",
				Description:   "testPartner",
				Latitude:      100,
				Longtitude:    100,
				Address:       "testPartner",
				City:          "testPartner",
			})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

			context := e.NewContext(req, res)
			context.SetPath("/partners/submission")

			partnerController := partner.NewPartnerController(mockPartnerRepository2{})
			if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.ApplyPartner())(context); err != nil {
				log.Fatal(err)
				return
			}

			var responses common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &responses)
			assert.Equal(t, "Successful Operation", responses.Message)
		})

	t.Run(
		"test Apply as Partner not found", func(t *testing.T) {
			e := echo.New()

			e.Validator = &partner.PartnerValidator{Validator: validator.New()}
			requestBody, _ := json.Marshal(partner.PartnerUserRequestFormat{
				BussinessName: "testPartner",
				Description:   "testPartner",
				Latitude:      100,
				Longtitude:    100,
				Address:       "testPartner",
				City:          "testPartner",
			})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

			context := e.NewContext(req, res)
			context.SetPath("/partners/submission")

			partnerController := partner.NewPartnerController(mockFalsePartnerRepository{})
			if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.ApplyPartner())(context); err != nil {
				log.Fatal(err)
				return
			}

			var responses common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &responses)
			assert.Equal(t, "Successful Operation", responses.Message)
		})
}

func TestGetAllPartner(t *testing.T) {
	t.Run(
		"test Get all Partner", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")

			context := e.NewContext(req, res)
			context.SetPath("/partners/submission")

			partnerController := partner.NewPartnerController(mockPartnerRepository{})
			partnerController.GetAllPartner()(context)

			var responses common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &responses)
			assert.Equal(t, "Successful Operation", responses.Message)
		})

	t.Run(
		"test Get all Partner failed", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")

			context := e.NewContext(req, res)
			context.SetPath("/partners/submission")

			partnerController := partner.NewPartnerController(mockFalsePartnerRepository{})
			partnerController.GetAllPartner()(context)

			var responses common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &responses)
			assert.Equal(t, "Bad Request", responses.Message)
		})

}

func TestAcceptPartner(t *testing.T) {
	t.Run("test accept partner", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController((mockPartnerRepository{}))
		partnerController.AcceptPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

	})

	t.Run("test accept partner bad request", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("a")

		partnerController := partner.NewPartnerController((mockPartnerRepository{}))
		partnerController.AcceptPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("test accept partner bad request 2", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController((mockPartnerRepository2{}))
		partnerController.AcceptPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("test accept partner bad request 3", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/accept")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController((mockFalsePartnerRepository{}))
		partnerController.AcceptPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})
}
func TestRejectPartner(t *testing.T) {
	t.Run("test reject partner", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController((mockPartnerRepository{}))
		partnerController.RejectPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

	})

	t.Run("test reject partner badrequest 1", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("a")

		partnerController := partner.NewPartnerController((mockPartnerRepository{}))
		partnerController.RejectPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("test reject partner badrequest 2", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController((mockPartnerRepository3{}))
		partnerController.RejectPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("test reject partner badrequest 3", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/:id/reject")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController((mockFalsePartnerRepository{}))
		partnerController.RejectPartner()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})
}

func TestGetPartner(t *testing.T) {
	t.Run("Test Get Partner Success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/partners/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		partnerController := partner.NewPartnerController(mockPartnerRepository{})
		partnerController.GetPartnerProduct()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})
	t.Run("Test Get Partner Failed", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/partners/:id")
		context.SetParamNames("id")
		context.SetParamValues("a")

		partnerController := partner.NewPartnerController(mockPartnerRepository{})
		partnerController.GetPartnerProduct()(context)

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})
}

func TestUpload(t *testing.T) {
	t.Run("login", func(t *testing.T) {

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
		// fmt.Println(response)
		JwtToken = response.Data.(string)
	})

	t.Run("test upload", func(t *testing.T) {

		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("legal_document", "test.pdf")

		file, _ := os.Open("test.pdf")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		// req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/partners/submission/upload")

		partnerController := partner.NewPartnerController(mockPartnerRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "", responses.Message)
	})

}

func TestReport(t *testing.T) {
	t.Run("login", func(t *testing.T) {

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
		// fmt.Println(response)
		JwtToken = response.Data.(string)
	})

	t.Run("test report", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/partners/report")

		partnerController := partner.NewPartnerController(mockPartnerRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.Report())(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

	})

	t.Run("test report failed", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/partners/report")

		partnerController := partner.NewPartnerController(mockFalsePartnerRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(partnerController.Report())(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)

	})

}

//======================
//MOCK PARTNER REPOSITORY
//======================
type mockPartnerRepository struct{}

func (m mockPartnerRepository) ApplyPartner(partner models.Partner) (models.Partner, error) {
	return models.Partner{
		Model:         gorm.Model{},
		UserID:        1,
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository) GetAllPartner() ([]models.Partner, error) {
	return []models.Partner{
		{
			UserID:        1,
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, nil
}

func (m mockPartnerRepository) GetPartner(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository) FindPartnerId(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository) FindUserId(userId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository) AcceptPartner(partner models.Partner) error {
	return nil
}

func (m mockPartnerRepository) RejectPartner(partner models.Partner) error {
	return nil
}

func (m mockPartnerRepository) GetAllPartnerProduct() ([]models.Partner, error) {
	return []models.Partner{
		{
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, nil
}

func (m mockPartnerRepository) UploadDocument(partnerID int, partner models.Partner) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository) Report(partnerId int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:         1,
			PartnerID:      1,
			Buffet:         false,
			Quantity:       1,
			DateTime:       time.Time{},
			Latitude:       1,
			Longtitude:     1,
			Distance:       1,
			TotalPrice:     1000,
			InvoiceID:      "1111",
			PaymentUrl:     "localhost",
			PaymentChannel: "BNI",
			PaymentMethod:  "BANK TRANSFER",
			PaidAt:         time.Time{},
			Status:         "PAID",
			Products: []models.Product{
				{

					Title: "udang",
					Type:  "snack",
					Price: 1000,
				},
			},
		},
	}, nil
}

//======================
//MOCK PARTNER REPOSITORY2
//======================
type mockPartnerRepository2 struct{}

func (m mockPartnerRepository2) ApplyPartner(partner models.Partner) (models.Partner, error) {
	return models.Partner{
		Model:         gorm.Model{},
		UserID:        1,
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		Status:        "reject",
	}, nil
}

func (m mockPartnerRepository2) GetAllPartner() ([]models.Partner, error) {
	return []models.Partner{
		{
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, nil
}

func (m mockPartnerRepository2) GetPartner(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository2) FindPartnerId(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "reject",
	}, nil
}

func (m mockPartnerRepository2) FindUserId(userId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "reject",
	}, nil
}

func (m mockPartnerRepository2) AcceptPartner(partner models.Partner) error {
	return nil
}

func (m mockPartnerRepository2) RejectPartner(partner models.Partner) error {
	return nil
}

func (m mockPartnerRepository2) GetAllPartnerProduct() ([]models.Partner, error) {
	return []models.Partner{
		{
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, nil
}

func (m mockPartnerRepository2) UploadDocument(partnerID int, partner models.Partner) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository2) Report(partnerId int) ([]models.Transaction, error) {
	return []models.Transaction{
		{
			UserID:         1,
			PartnerID:      1,
			Buffet:         false,
			Quantity:       1,
			DateTime:       time.Time{},
			Latitude:       1,
			Longtitude:     1,
			Distance:       1,
			TotalPrice:     1000,
			InvoiceID:      "1111",
			PaymentUrl:     "localhost",
			PaymentChannel: "BNI",
			PaymentMethod:  "BANK TRANSFER",
			PaidAt:         time.Time{},
			Status:         "PAID",
		},
	}, nil
}

//======================
//MOCK PARTNER REPOSITORY3
//======================
type mockPartnerRepository3 struct{}

func (m mockPartnerRepository3) ApplyPartner(partner models.Partner) (models.Partner, error) {
	return models.Partner{
		Model:         gorm.Model{},
		UserID:        1,
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		Status:        "reject",
	}, nil
}

func (m mockPartnerRepository3) GetAllPartner() ([]models.Partner, error) {
	return []models.Partner{
		{
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, nil
}

func (m mockPartnerRepository3) GetPartner(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository3) FindPartnerId(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "active",
	}, nil
}

func (m mockPartnerRepository3) FindUserId(userId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "reject",
	}, nil
}

func (m mockPartnerRepository3) AcceptPartner(partner models.Partner) error {
	return nil
}

func (m mockPartnerRepository3) RejectPartner(partner models.Partner) error {
	return nil
}

func (m mockPartnerRepository3) GetAllPartnerProduct() ([]models.Partner, error) {
	return []models.Partner{
		{
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, nil
}

func (m mockPartnerRepository3) UploadDocument(partnerID int, partner models.Partner) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, nil
}

func (m mockPartnerRepository3) Report(partnerId int) ([]models.Transaction, error) {

	return []models.Transaction{
		{
			UserID:         1,
			PartnerID:      1,
			Buffet:         false,
			Quantity:       1,
			DateTime:       time.Time{},
			Latitude:       1,
			Longtitude:     1,
			Distance:       1,
			TotalPrice:     1000,
			InvoiceID:      "1111",
			PaymentUrl:     "localhost",
			PaymentChannel: "BNI",
			PaymentMethod:  "BANK TRANSFER",
			PaidAt:         time.Time{},
			Status:         "PAID",
			Products: []models.Product{
				{
					Title: "rendang",
				},
				{
					Title: "rendang2",
				},
			},
		},
	}, nil
}

//======================
//MOCK FALSE PARTNER  REPOSITORY
//======================
type mockFalsePartnerRepository struct{}

func (m mockFalsePartnerRepository) ApplyPartner(partner models.Partner) (models.Partner, error) {
	return models.Partner{
		Model:         gorm.Model{},
		UserID:        1,
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		Status:        "DRAFT",
	}, errors.New("failed")
}

func (m mockFalsePartnerRepository) GetAllPartner() ([]models.Partner, error) {
	return nil, errors.New("failed")
}

func (m mockFalsePartnerRepository) GetPartner(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, errors.New("failed")
}

func (m mockFalsePartnerRepository) FindPartnerId(partnerId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, errors.New("failed")
}

func (m mockFalsePartnerRepository) FindUserId(userId int) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, errors.New("failed")
}

func (m mockFalsePartnerRepository) AcceptPartner(partner models.Partner) error {
	return errors.New("failed")
}

func (m mockFalsePartnerRepository) RejectPartner(partner models.Partner) error {
	return errors.New("failed")
}

func (m mockFalsePartnerRepository) GetAllPartnerProduct() ([]models.Partner, error) {
	return []models.Partner{
		{
			BussinessName: "testPartner",
			Description:   "testPartner",
			Latitude:      100,
			Longtitude:    100,
			Address:       "testPartner",
			City:          "testPartner",
			LegalDocument: "testPartner.pdf",
			Status:        "DRAFT",
		},
	}, errors.New("failed")
}

func (m mockFalsePartnerRepository) UploadDocument(partnerID int, partner models.Partner) (models.Partner, error) {
	return models.Partner{
		BussinessName: "testPartner",
		Description:   "testPartner",
		Latitude:      100,
		Longtitude:    100,
		Address:       "testPartner",
		City:          "testPartner",
		LegalDocument: "testPartner.pdf",
		Status:        "DRAFT",
	}, errors.New("failed")
}

func (m mockFalsePartnerRepository) Report(partnerId int) ([]models.Transaction, error) {
	return nil, errors.New("failed")
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
