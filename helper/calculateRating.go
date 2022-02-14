package helper

import (
	"github.com/furqonzt99/snackbox/models"
)

func CalculateRating(ratings []models.Rating) (avgRating float64) {
	if len(ratings) < 1 {
		return 0
	}
	
	for _, rating := range ratings {
		avgRating += float64(rating.Rating)
	}

	return avgRating / float64(len(ratings))
}