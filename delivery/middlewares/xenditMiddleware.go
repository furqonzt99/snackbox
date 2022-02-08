package middlewares

import (
	"net/http"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/labstack/echo/v4"
)

func CheckXHeaderToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := c.Request()
		headers := req.Header

		xCallbackToken := headers.Get("X-Callback-Token")

		if xCallbackToken != constants.XENDIT_CALLBACK_TOKEN {
			return c.JSON(http.StatusUnauthorized, common.NewUnauthorizeResponse())
		}

		return next(c)
	}
}