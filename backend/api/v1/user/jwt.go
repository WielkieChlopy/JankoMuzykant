package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func generateJWT(id uuid.UUID, username string, secretKey []byte) (string, error) {
	claims := jwtCustomClaims{
		id,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func userIDFromToken(c echo.Context) (uuid.UUID, error) {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	return claims.Id, nil
}
