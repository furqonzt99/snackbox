package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterProductPath(e *echo.Echo, productCtrl *product.ProductController) {

	e.POST("/products", productCtrl.AddProduct(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))

}
