package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/rating"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRatingPath(e *echo.Echo, RatingController *rating.RatingController) {

	e.POST("/ratings", RatingController.Create, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckUserRole)
}