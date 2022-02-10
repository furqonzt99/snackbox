package middlewares

import (
	"net/http"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/labstack/echo/v4"
)

func CheckAdminRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := ExtractTokenUser(c)

		if user.Role != "admin" {
			return c.JSON(http.StatusUnauthorized, common.NewUnauthorizeResponse())
		}
		return next(c)
	}
}