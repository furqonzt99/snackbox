package middlewares

import (
	"errors"
	"time"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CreateToken(userId, partnerId int, email, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = int(userId)
	claims["partnerId"] = int(partnerId)
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constants.JWT_SECRET_KEY))
}

func ExtractTokenUser(e echo.Context) (common.JWTPayload, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userId"].(float64)
		partnerId := claims["partnerId"].(float64)
		email := claims["email"]
		role := claims["role"]
		return common.JWTPayload{
			UserID: int(userId),
			PartnerID: int(partnerId),
			Email:  email.(string),
			Role:  role.(string),
		}, nil
	}
	return common.JWTPayload{}, errors.New("invalid token")
}