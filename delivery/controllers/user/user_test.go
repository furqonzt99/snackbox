package user_test

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

func TestRegisterUser(t *testing.T) {
	//======================
	//REGISTER TEST
	//======================
	t.Run("Test Register", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
			"name":     "tester",
			"city":     "cityTest",
			"address":  "addressTest",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/register")

		userController := user.NewUsersControllers(mockUserRepository{})
		userController.RegisterController()(context)

		response := user.RegisterUserResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Register Password Length Below 8", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test",
			"name":     "tester",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/register")

		userController := user.NewUsersControllers(mockFalseUserRepository{})
		userController.RegisterController()(context)

		response := user.RegisterUserResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Email Already Exist", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
			"name":     "tester",
			"city":     "cityTest",
			"address":  "addressTest",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/register")

		userController := user.NewUsersControllers(mockFalseUserRepository{})
		userController.RegisterController()(context)

		response := user.RegisterUserResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Email already exist", response.Message)
	})
}

var JwtToken string //TOKEN FROM LOGIN

func TestLoginUser(t *testing.T) {
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

	t.Run("Error Test Login Password Length Below 4", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "tes",
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

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Login not found", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test22@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := user.NewUsersControllers(mockFalseUserRepository{})
		userController.LoginController()(context)
		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "User not found", response.Message)
	})

	t.Run("Error Test Login Wrong Password", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "1test1234",
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

		assert.Equal(t, "Wrong Password", response.Message)
	})

	t.Run("Test Login err", func(t *testing.T) {
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

		userController := user.NewUsersControllers(mockUserRepository2{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		JwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, JwtToken)
	})

	//======================
	//TEST GET USER PROFILE
	//======================
	t.Run("Test Get User", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		// fmt.Println(jwtToken)
		context := e.NewContext(req, res)
		context.SetPath("/profile")

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.GetUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		// assert.Equal(t, response.Data.(map[string]interface{})["name"], "tester")
	})

	t.Run("Test Get User partner", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))
		// fmt.Println(jwtToken)
		context := e.NewContext(req, res)
		context.SetPath("/profile")

		userController := user.NewUsersControllers(mockUserRepository2{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.GetUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		// assert.Equal(t, response.Data.(map[string]interface{})["name"], "tester")
	})
	//======================
	//TEST UPDATE USER PROFILE
	//======================
	t.Run("Test Update", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test2@gmail.com",
			"password": "test4321",
			"name":     "tester2",
			"address":  "alamat",
			"city":     "kota",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.UpdateUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Update Password Length Below 4", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test2@gmail.com",
			"password": "tes",
			"name":     "tester2",
			"address":  "alamat",
			"city":     "kota",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := user.NewUsersControllers(mockFalseUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.UpdateUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Update User Repo", func(t *testing.T) {
		e := echo.New()
		e.Validator = &user.UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test2@gmail.com",
			"password": "test1234",
			"name":     "tester2",
			"address":  "alamat",
			"city":     "kota",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := user.NewUsersControllers(mockFalseUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.UpdateUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	//======================
	//TEST DELETE USER PROFILE
	//======================
	t.Run("Test Delete", func(t *testing.T) {
		e := echo.New()
		// mw.LogMiddleware(e)
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", JwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.DeleteUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
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
		fw, _ := writer.CreateFormFile("photo", "golang.jpeg") // add file to partner folder

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

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("test upload err repo.Get", func(t *testing.T) {

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

		userController := user.NewUsersControllers(mockFalseUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "False Login Object", responses.Message)
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
		fw, _ := writer.CreateFormFile("photo", "golang.jpeg") // add file to partner folder

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

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

	t.Run("test upload wrong extention", func(t *testing.T) {

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
		fw, _ := writer.CreateFormFile("photo", "aku.txt") // add file to partner folder

		file, _ := os.Open("aku.txt")

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

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "file type must an image", responses.Message)
	})

	t.Run("test upload err file photo", func(t *testing.T) {

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
		fw, _ := writer.CreateFormFile("photo1", "golang.jpeg") // add file to partner folder

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

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "http: no such file", responses.Message)
	})

	t.Run("test upload success ", func(t *testing.T) {

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
		fw, _ := writer.CreateFormFile("photo", "golang.jpeg") // add file to partner folder

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

		userController := user.NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constants.JWT_SECRET_KEY))(userController.Upload)(context); err != nil {
			log.Fatal(err)
			return
		}

		var responses common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &responses)
		assert.Equal(t, "Successful Operation", responses.Message)
	})

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
		Name:     "tester",
		Photo:    "adaaja"}, nil
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

//======================
//MOCK USER REPOSITORY2
//======================
type mockUserRepository2 struct{}

func (m mockUserRepository2) Register(newUser models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	return models.User{
		Email:    newUser.Email,
		Password: string(hash),
		Name:     newUser.Name,
		Address:  newUser.Address,
		City:     newUser.City,
	}, nil
}

func (m mockUserRepository2) Login(email string) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return models.User{
		Email:    "test@gmail.com",
		Password: string(hash),
		Partner: models.Partner{
			Model: gorm.Model{
				ID: 2,
			},
			Status: "active",
		},
	}, nil
}

func (m mockUserRepository2) Get(userid int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return models.User{
		Email:    "test@gmail.com",
		Password: string(hash),
		Name:     "tester",
		Photo:    "poto",
		Partner: models.Partner{
			Model: gorm.Model{
				ID: 1,
			},
		},
	}, nil
}

func (m mockUserRepository2) Update(newUser models.User, userId int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{
		Email:    "test2@gmail.com",
		Password: string(hash),
		Name:     "tester2",
		Photo:    "apaja",
	}, nil
}

func (m mockUserRepository2) Delete(userId int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{
		Email:    "test2@gmail.com",
		Password: string(hash), Name: "tester2",
	}, nil
}

//======================
//MOCK FALSE REPOSITORY
//======================
type mockFalseUserRepository struct{}

func (m mockFalseUserRepository) Register(newUser models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	return models.User{Email: newUser.Email, Password: string(hash), Name: newUser.Name}, errors.New("Email already exist")
}

func (m mockFalseUserRepository) Login(email string) (models.User, error) {
	// hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{}, errors.New("email not found")
}

func (m mockFalseUserRepository) Get(userid int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return models.User{Email: "test@gmail.com", Password: string(hash), Name: "tester"}, errors.New("False Login Object")
}

func (m mockFalseUserRepository) Update(newUser models.User, userId int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{Email: "test2@gmail.com", Password: string(hash), Name: "tester2"}, errors.New("False Login Object")
}

func (m mockFalseUserRepository) Delete(userId int) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return models.User{Email: "test2@gmail.com", Password: string(hash), Name: "tester2"}, errors.New("False Login Object")
}
