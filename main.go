package main

import (
	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/delivery/routes"
	pt "github.com/furqonzt99/snackbox/repositories/partner"
	ur "github.com/furqonzt99/snackbox/repositories/user"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.GetConfig()

	db := utils.InitDB(config)

	utils.InitialMigrate(db)

	//repo
	userRepo := ur.NewUserRepo(db)
	partnerRepo := pt.NewPartnerRepo(db)
	//controller
	userCtrl := user.NewUsersControllers(userRepo)
	partnerCtrl := partner.NewPartnerController(partnerRepo)
	e := echo.New()
	middlewares.LogMiddleware(e)
	e.Pre(middleware.RemoveTrailingSlash())

	//validator
	e.Validator = &user.UserValidator{Validator: validator.New()}
	e.Validator = &partner.PartnerValidator{Validator: validator.New()}

	//routes
	routes.RegisterUserPath(e, userCtrl)
	routes.RegisterPartnerPath(e, partnerCtrl)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
