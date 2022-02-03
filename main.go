package main

import (
	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.GetConfig()

	db := utils.InitDB(config)

	utils.InitialMigrate(db)

	e := echo.New()
	
	middlewares.LogMiddleware(e)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Logger.Fatal(e.Start(":" + config.Port))
}