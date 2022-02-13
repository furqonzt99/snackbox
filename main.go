package main

import (
	config "github.com/furqonzt99/snackbox/configs"
	"github.com/furqonzt99/snackbox/delivery/controllers/bank"
	"github.com/furqonzt99/snackbox/delivery/controllers/cashout"
	"github.com/furqonzt99/snackbox/delivery/controllers/partner"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/controllers/rating"
	"github.com/furqonzt99/snackbox/delivery/controllers/transaction"
	"github.com/furqonzt99/snackbox/delivery/controllers/user"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/delivery/routes"
	br "github.com/furqonzt99/snackbox/repositories/bank"
	cr "github.com/furqonzt99/snackbox/repositories/cashout"
	pt "github.com/furqonzt99/snackbox/repositories/partner"
	pd "github.com/furqonzt99/snackbox/repositories/product"
	rr "github.com/furqonzt99/snackbox/repositories/rating"
	tr "github.com/furqonzt99/snackbox/repositories/transaction"
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
	productRepo := pd.NewProductRepo(db)
	transactionRepo := tr.NewTransactionRepository(db)
	ratingRepo := rr.NewRatingRepository(db)
	cashoutRepo := cr.NewCashoutRepository(db)
	bankRepo := br.NewBankRepository(db)

	//controller
	userCtrl := user.NewUsersControllers(userRepo)
	partnerCtrl := partner.NewPartnerController(partnerRepo)
	productCtrl := product.NewProductController(productRepo)
	transactionController := transaction.NewTransactionController(transactionRepo)
	ratingController := rating.NewRatingController(ratingRepo)
	cashoutController := cashout.NewCashoutController(cashoutRepo)
	bankController := bank.NewBankController(bankRepo)

	//echo package
	e := echo.New()
	middlewares.LogMiddleware(e)
	e.Pre(middleware.RemoveTrailingSlash())

	//validator
	e.Validator = &user.UserValidator{Validator: validator.New()}
	e.Validator = &partner.PartnerValidator{Validator: validator.New()}
	e.Validator = &product.ProductValidator{Validator: validator.New()}
	e.Validator = &transaction.TransactionValidator{Validator: validator.New()}
	e.Validator = &rating.RatingValidator{Validator: validator.New()}
	e.Validator = &cashout.CashoutValidator{Validator: validator.New()}

	//routes
	routes.RegisterUserPath(e, userCtrl)
	routes.RegisterPartnerPath(e, partnerCtrl)
	routes.RegisterProductPath(e, productCtrl)
	routes.RegisterTransactionPath(e, transactionController)
	routes.RegisterRatingPath(e, ratingController)
	routes.RegisterCashoutPath(e, cashoutController)
	routes.RegisterBankPath(e, bankController)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
