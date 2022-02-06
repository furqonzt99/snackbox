package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckUserRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := ExtractTokenUser(c)

		if user.Role != "user" {
			return c.JSON(
				http.StatusUnauthorized, map[string]interface{}{
					"Message": "Unauthorized",
				},
			)
		}
		return next(c)
	}
}