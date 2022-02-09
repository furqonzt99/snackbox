package user

import (
	"net/http"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/user"
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
			Address:  newUserReq.Address,
			City:     newUserReq.City,
			Password: string(hash),
		}

		res, err := uscon.Repo.Register(newUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(406, "Email already exist"))
		}

		data := UserResponse{
			ID:      res.ID,
			Name:    res.Name,
			Email:   res.Email,
			Address: res.Address,
			City:    res.City,
			Balance: res.Balance,
			Role:    res.Role,
		}
		return c.JSON(http.StatusOK, common.SuccessResponse(data))
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
			return echo.NewHTTPError(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}

		hash, err := helper.Checkpwd(user.Password, loginUser.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(403, "Wrong Password"))
		}

		var token string

		var partnerId int
		if user.Partner.ID != 0 {
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

		data := UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Address: user.Address,
			City:    user.City,
			Balance: user.Balance,
			Role:    user.Role,
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
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
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
