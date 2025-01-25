package songs

import (
	"github.com/labstack/echo/v4"
)

type songRequest struct {
	Url string `json:"url" validate:"required"`
}

func (r *songRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}
