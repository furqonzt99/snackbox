package product_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JwtToken string

func TestProduct(t *testing.T) {
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

}

//======================
//MOCK PRODUCT REPOSITORY
//======================

type mockProductRepository struct{}

func (m mockProductRepository) AddProduct(product models.Product) (models.Product, error) {
	return models.Product{
		// PartnerID:   0
		Title:       "testProduct1",
		Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, nil
}

func (m mockProductRepository) FindProduct(productId, partnerId int) (models.Product, error) {
	return models.Product{
		Model: gorm.Model{
			ID: uint(productId),
		},
		PartnerID:   uint(partnerId),
		Title:       "testProduct1",
		Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, nil
}
func (m mockProductRepository) DeleteProduct(productId, partnerId int) error {
	return nil
}
func (m mockProductRepository) GetAllProduct(offset, pageSize int, search string) ([]models.Product, error) {
	return []models.Product{
		{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		},
	}, nil
}

func (m mockProductRepository) UploadImage(productID int, product models.Product) (models.Product, error) {
	return models.Product{
		Model: gorm.Model{
			ID: uint(productID),
		},
		// PartnerID:   0,
		Title: "testProduct1",
		// Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, nil
}

//======================
//MOCK FALSE PRODUCT REPOSITORY
//======================

type mockFalseProductRepository struct{}

func (m mockFalseProductRepository) AddProduct(product models.Product) (models.Product, error) {
	return models.Product{
		// PartnerID:   0
		Title:       "testProduct1",
		Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, errors.New("failed")
}

func (m mockFalseProductRepository) FindProduct(productId, partnerId int) (models.Product, error) {
	return models.Product{
		Model: gorm.Model{
			ID: uint(productId),
		},
		PartnerID:   uint(partnerId),
		Title:       "testProduct1",
		Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, errors.New("failed")
}
func (m mockFalseProductRepository) DeleteProduct(productId, partnerId int) error {
	return errors.New("failed")
}
func (m mockFalseProductRepository) GetAllProduct(offset, pageSize int, search string) ([]models.Product, error) {
	return nil, errors.New("failed")
}

func (m mockFalseProductRepository) UploadImage(productID int, product models.Product) (models.Product, error) {
	return models.Product{
		Model: gorm.Model{
			ID: uint(productID),
		},
		// PartnerID:   0,
		Title: "testProduct1",
		// Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, errors.New("failed")
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
