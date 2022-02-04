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
		var login UserLoginRequestFormat
		c.Bind(&login)

		if err := c.Validate(login); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		user, err := uscon.Repo.Login(login.Email)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}

		hash, err := helper.Checkpwd(user.Password, login.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(403, "Wrong Password"))
		}

		var token string

		if hash {
			token, _ = middlewares.CreateToken(int(user.ID), user.Email, user.Role)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(token))
	}
}

func (uscon UserController) GetUserController() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		user, _ := uscon.Repo.Get(userJwt.UserID)

		data := UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (uscon UserController) UpdateUserController() echo.HandlerFunc {

	return func(c echo.Context) error {
		user, _ := middlewares.ExtractTokenUser(c)

		updateUserReq := PutUserRequestFormat{}
		c.Bind(&updateUserReq)

		if err := c.Validate(updateUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		updateUser := models.User{}
		updateUser.Email = updateUserReq.Email
		updateUser.Name = updateUserReq.Name

		if updateUserReq.Password != "" {
			hash, _ := bcrypt.GenerateFromPassword([]byte(updateUserReq.Password), 14)
			updateUser.Password = string(hash)
		}

		userData, err := uscon.Repo.Update(updateUser, user.UserID)

		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}

		data := UserResponse{
			ID:    uint(user.UserID),
			Name:  userData.Name,
			Email: userData.Email,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}
}

func (uscon UserController) DeleteUserController() echo.HandlerFunc {

	return func(c echo.Context) error {
		userId, _ := middlewares.ExtractTokenUser(c)

		uscon.Repo.Delete(userId.UserID)

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}
