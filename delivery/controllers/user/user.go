package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/user"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Repo user.UserInterface
}

func NewUsersControllers(usrep user.UserInterface) *UserController {
	return &UserController{Repo: usrep}
}

func (uscon UserController) RegisterController() echo.HandlerFunc {
	return func(c echo.Context) error {
		newUserReq := RegisterUserRequestFormat{}
		c.Bind(&newUserReq)

		if err := c.Validate(newUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(newUserReq.Password), 14)
		newUser := models.User{
			Name:     newUserReq.Name,
			Email:    newUserReq.Email,
			Password: string(hash),
			City:     newUserReq.City,
			Address:  newUserReq.Address,
		}

		_, err := uscon.Repo.Register(newUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(406, "Email already exist"))
		}

		// data := UserResponse{
		// 	ID:      res.ID,
		// 	Name:    res.Name,
		// 	Email:   res.Email,
		// 	Address: res.Address,
		// 	City:    res.City,
		// 	Balance: res.Balance,
		// 	Role:    res.Role,
		// }
		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (uscon UserController) LoginController() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginUser UserLoginRequestFormat
		c.Bind(&loginUser)

		if err := c.Validate(loginUser); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		user, err := uscon.Repo.Login(loginUser.Email)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}

		hash, err := helper.Checkpwd(user.Password, loginUser.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Wrong Password"))
		}

		var token string

		var partnerId int
		const STATUS_ACTIVE = "active"
		if user.Partner.ID != 0 && user.Partner.Status == STATUS_ACTIVE {
			partnerId = int(user.Partner.ID)
		}

		if hash {
			token, _ = middlewares.CreateToken(int(user.ID), partnerId, user.Email, user.Role)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(token))
	}
}

func (uscon UserController) GetUserController() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		user, _ := uscon.Repo.Get(userJwt.UserID)

		var userProfile string
		if user.Photo != "" { 
			userProfile = fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, user.Photo)
		}

		if user.Partner.ID == 0 {
			data := UserProfileResponse{
				ID:      user.ID,
				Name:    user.Name,
				Photo: userProfile,
				Email:   user.Email,
				Balance: user.Balance,
			}
			return c.JSON(http.StatusOK, common.SuccessResponse(data))
		}

		data := UserProfileResponseWithPartner{
			ID:      user.ID,
			Name:    user.Name,
			Photo: userProfile,
			Email:   user.Email,
			Balance: user.Balance,
			Partner: partner.GetPartnerProfileResponse{
				ID:            int(user.Partner.ID),
				BussinessName: user.Partner.BussinessName,
				Description:   user.Partner.Description,
				Latitude:      user.Partner.Latitude,
				Longtitude:    user.Partner.Longtitude,
				LegalDocument: user.Partner.LegalDocument,
				Status:        user.Partner.Status,
			},
		}
		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (uscon UserController) UpdateUserController() echo.HandlerFunc {

	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		updateUserReq := PutUserRequestFormat{}
		c.Bind(&updateUserReq)

		if err := c.Validate(updateUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		updateUser := models.User{}
		updateUser.Email = updateUserReq.Email
		updateUser.Name = updateUserReq.Name
		updateUser.Address = updateUserReq.Address
		updateUser.City = updateUserReq.City

		if updateUserReq.Password != "" {
			hash, _ := bcrypt.GenerateFromPassword([]byte(updateUserReq.Password), 14)
			updateUser.Password = string(hash)
		}

		_, err := uscon.Repo.Update(updateUser, userJwt.UserID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (uscon UserController) DeleteUserController() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, _ := middlewares.ExtractTokenUser(c)

		uscon.Repo.Delete(userId.UserID)

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (uc UserController) Upload(c echo.Context) error {
	var requestUpload UserPhotoRequest

	if err := c.Bind(&requestUpload); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userDB, err := uc.Repo.Get(user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.ErrorResponse(http.StatusNotFound, err.Error()))
	}
	
	file, err := c.FormFile("photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	defer src.Close()

	head := make([]byte, 261)
  	src.Read(head)

	kind, _ := filetype.Match(head)

	if !filetype.IsImage(head) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "file type must an image"))
	}

	prefix := "profiles/"

	fileID := strings.ReplaceAll(uuid.New().String(), "-", "")
	file.Filename = fmt.Sprint(prefix, fileID, ".", kind.Extension)

	if userDB.Photo != "" {
		if err := helper.GetObjectS3(userDB.Photo); err == nil {
			_ = helper.DeleteObjectS3(userDB.Photo)
		}
	}

	if err := helper.UploadObjectS3(file.Filename, src); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userProfile := models.User{
		Photo: file.Filename,
	}

	_, err = uc.Repo.Update(userProfile, int(userDB.ID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}
