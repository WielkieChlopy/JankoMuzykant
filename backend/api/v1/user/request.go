package user

import (
	"backend/models"

	"github.com/labstack/echo/v4"
)

type userUpdateRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userUpdateRequest) bind(c echo.Context) (*models.User, error) {
	u := &models.User{}
	if err := c.Bind(r); err != nil {
		return nil, err
	}

	if err := c.Validate(r); err != nil {
		return nil, err
	}

	u.Username = r.User.Username

	h, err := HashPassword(r.User.Password)
	if err != nil {
		return nil, err
	}

	u.Password = h

	return u, nil
}

type userLoginRequest struct {
	User struct {
		Username string `json:"Username" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}
