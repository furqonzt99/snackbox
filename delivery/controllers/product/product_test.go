package product_test

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

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/models"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JwtToken string

func TestAddProduct(t *testing.T) {
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
		JwtToken = response.Data.(string)

	})

	t.Run("add product Success", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.AddProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

	})
	t.Run("add product bad reqeust 1", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title: "testProduct1",
			Type:  "testProduct1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.AddProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("add product bad reqeust 2", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products")

		userController := product.NewProductController(mockFalseProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.AddProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})
}
func TestPutProduct(t *testing.T) {
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
		JwtToken = response.Data.(string)

	})

	t.Run("Put product Success", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.PutProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

	})

	t.Run("Put product badrequest 1", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("a")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.PutProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("Put product badrequest 2", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Type:  "testProduct1",
			Price: 1000,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.PutProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("Put product badrequest 3", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := product.NewProductController(mockFalseProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.PutProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)

	})
	t.Run("Put product badrequest 4", func(t *testing.T) {

		e := echo.New()
		e.Validator = &product.ProductValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(product.RegisterProductRequestFormat{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := product.NewProductController(mockFalseProductRepository2{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.PutProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})
}

func TestDeleteProduct(t *testing.T) {
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
		JwtToken = response.Data.(string)

	})

	t.Run("delete product Success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.DeleteProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

	})

	t.Run("delete product badrequest 1", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("a")

		userController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.DeleteProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)

	})

	t.Run("delete product badrequest 2", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := product.NewProductController(mockFalseProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.DeleteProduct())(context); err != nil {
			log.Fatal(err)
			return
		}

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)

	})
}

func TestGetAllProduct(t *testing.T) {
	t.Run("get all product Success", func(t *testing.T) {

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("search")
		context.SetParamValues("product")

		userController := product.NewProductController(mockProductRepository{})
		userController.GetAllProduct()(context)

		responses := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)

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
		JwtToken = response.Data.(string)

	})
	t.Run("test upload", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image", "golang.jpeg") // add file to partner folder

		file, _ := os.Open("golang.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		productController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("test upload bad request param", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image", "golang.jpeg") // add file to partner folder

		file, _ := os.Open("golang.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("a")

		productController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
	})

	t.Run("test upload err find product", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image", "golang.jpeg") // add file to partner folder

		file, _ := os.Open("golang.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		productController := product.NewProductController(mockFalseProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Not Found", responses.Message)
	})
	t.Run("test upload err form file", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image_failed", "golang.jpeg") // add file to partner folder

		file, _ := os.Open("golang.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		productController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "http: no such file", responses.Message)
	})

	t.Run("test upload err extension not image", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image", "golang1.jpeg") // add file to partner folder

		file, _ := os.Open("golang1.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		productController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "file type must an image", responses.Message)
	})

	t.Run("test upload err upload s3", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION_Failed")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image", "golang.jpeg") // add file to partner folder

		file, _ := os.Open("golang.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		productController := product.NewProductController(mockProductRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "MissingRegion: could not find region configuration", responses.Message)
	})

	t.Run("test upload err repo.upload image", func(t *testing.T) {

		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		constants.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
		constants.AWS_ACCESS_SECRET_KEY = os.Getenv("AWS_ACCESS_SECRET_KEY")
		constants.S3_REGION = os.Getenv("S3_REGION")
		constants.S3_BUCKET = os.Getenv("S3_BUCKET")
		constants.LINK_TEMPLATE = os.Getenv("LINK_TEMPLATE")
		//////////////////////////////////
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		fw, _ := writer.CreateFormFile("image", "golang.jpeg") // add file to partner folder

		file, _ := os.Open("golang.jpeg")

		_, _ = io.Copy(fw, file)

		writer.Close()

		//////////////////////////////////
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", writer.FormDataContentType()) // <<< important part
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/products/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		productController := product.NewProductController(mockProductRepository5{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(productController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Bad Request", responses.Message)
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
func (m mockProductRepository) GetAllProduct(offset, pageSize int, search, category string, latitude, longtitude float64) ([]models.Product, error) {
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
//MOCK PRODUCT REPOSITORY 5
//======================

type mockProductRepository5 struct{}

func (m mockProductRepository5) AddProduct(product models.Product) (models.Product, error) {
	return models.Product{
		// PartnerID:   0
		Title:       "testProduct1",
		Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, nil
}

func (m mockProductRepository5) FindProduct(productId, partnerId int) (models.Product, error) {
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
func (m mockProductRepository5) DeleteProduct(productId, partnerId int) error {
	return nil
}
func (m mockProductRepository5) GetAllProduct(offset, pageSize int, search, category string, latitude, longtitude float64) ([]models.Product, error) {
	return []models.Product{
		{
			Title:       "testProduct1",
			Type:        "testProduct1",
			Description: "testProduct1",
			Price:       1000,
		},
	}, nil
}

func (m mockProductRepository5) UploadImage(productID int, product models.Product) (models.Product, error) {
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
	}, errors.New("FAILED")
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
func (m mockFalseProductRepository) GetAllProduct(offset, pageSize int, search, category string, latitude, longtitude float64) ([]models.Product, error) {
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
//MOCK FALSE PRODUCT REPOSITORY2
//======================

type mockFalseProductRepository2 struct{}

func (m mockFalseProductRepository2) AddProduct(product models.Product) (models.Product, error) {
	return models.Product{
		// PartnerID:   0
		Title:       "testProduct1",
		Image:       "",
		Type:        "testProduct1",
		Description: "testProduct1",
		Price:       1000,
	}, errors.New("failed")
}

func (m mockFalseProductRepository2) FindProduct(productId, partnerId int) (models.Product, error) {
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
func (m mockFalseProductRepository2) DeleteProduct(productId, partnerId int) error {
	return errors.New("failed")
}
func (m mockFalseProductRepository2) GetAllProduct(offset, pageSize int, search, category string, latitude, longtitude float64) ([]models.Product, error) {
	return nil, errors.New("failed")
}

func (m mockFalseProductRepository2) UploadImage(productID int, product models.Product) (models.Product, error) {
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
