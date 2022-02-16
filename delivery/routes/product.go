package routes

import (
	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterProductPath(e *echo.Echo, productCtrl *product.ProductController) {

	e.POST("/products", productCtrl.AddProduct(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
	e.PUT("/products/:id", productCtrl.PutProduct(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
	e.DELETE("/products/:id", productCtrl.DeleteProduct(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
	e.GET("/products", productCtrl.GetAllProduct(), middleware.JWT([]byte(constants.JWT_SECRET_KEY)))
	e.PUT("/products/:id/image", productCtrl.Upload, middleware.JWT([]byte(constants.JWT_SECRET_KEY)), middlewares.CheckPartnerRole)
}
