package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterUserPath(e *echo.Echo, userCtrl *user.UserController) {

	e.POST("/register", userCtrl.RegisterController())
	e.POST("/login", userCtrl.LoginController())
	e.GET("/user", userCtrl.GetUserController(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.PUT("/user", userCtrl.UpdateUserController(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.DELETE("/user", userCtrl.DeleteUserController(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.PUT("/user/photo", userCtrl.Upload, middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
}
