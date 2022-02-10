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
	e.GET("/profile", userCtrl.GetUserController(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.PUT("/users", userCtrl.UpdateUserController(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.DELETE("/users", userCtrl.DeleteUserController(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
}
