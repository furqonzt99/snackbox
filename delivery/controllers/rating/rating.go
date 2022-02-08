package rating

import (
	"net/http"

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
		PartnerID: uint(ratingRequest.PartnerID),
		UserID:  uint(user.UserID),
		Rating:  ratingRequest.Rating,
		Comment: ratingRequest.Comment,
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