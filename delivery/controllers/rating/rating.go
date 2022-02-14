package rating

import (
	"net/http"
	"strconv"

	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/rating"
	"github.com/labstack/echo/v4"
)

type RatingController struct {
	Repo rating.RatingInterface
}

func NewRatingController(repo rating.RatingInterface) *RatingController {
	return &RatingController{Repo: repo}
}

func (rc RatingController) Create(c echo.Context) error {
	var ratingRequest PostRatingRequest

	if err := c.Bind(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	if err := c.Validate(&ratingRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, _ := middlewares.ExtractTokenUser(c)

	isCanGiveRating, _ := rc.Repo.IsCanGiveRating(user.UserID, ratingRequest.PartnerID)
	if !isCanGiveRating {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	data := models.Rating{
		TransactionID: 0,
		PartnerID:     uint(ratingRequest.PartnerID),
		UserID:        uint(user.UserID),
		Rating:        ratingRequest.Rating,
		Comment:       ratingRequest.Comment,
		Transaction:   models.Transaction{},
		User:          models.User{},
		Partner:       models.Partner{},
	}

	ratingData, err := rc.Repo.Create(data)
	if err != nil {
		ratingData, err = rc.Repo.Update(data)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
	}

	response := RatingResponse{
		PartnerID:  int(ratingData.PartnerID),
		UserID:   int(ratingData.UserID),
		Username: ratingData.User.Name,
		Rating:   ratingData.Rating,
		Comment:  ratingData.Comment,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response)) 
}

func (rc RatingController) GetByTrxID(c echo.Context) error {
	trxID, err := strconv.Atoi(c.Param("trxID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	rating, err := rc.Repo.GetByTrxID(trxID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	response := RatingResponse{
		PartnerID: int(rating.PartnerID),
		UserID:    int(rating.UserID),
		Username:  rating.User.Name,
		Rating:    rating.Rating,
		Comment:   rating.Comment,
	}

	return c.JSON(http.StatusBadRequest, common.SuccessResponse(response))

}