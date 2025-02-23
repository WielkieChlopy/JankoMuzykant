package queue

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type addSongRequest struct {
	URL string `json:"url" validate:"required"`
}

func (r *addSongRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type removeSongRequest struct {
	SongID uuid.UUID `json:"song_id" validate:"required"`
}

func (r *removeSongRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}
