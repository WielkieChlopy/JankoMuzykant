package playlist

import (
	"github.com/labstack/echo/v4"
)

type addSongRequest struct {
	URL string `json:"url" validate:"required"`
	//PlaylistID uuid.UUID `json:"playlist_id" validate:"required"`
}

type createPlaylistRequest struct {
	Name string `json:"name" validate:"required"`
}

type editPlaylistRequest struct {
	Name string `json:"name" validate:"required"`
}

type reorderPlaylistRequest struct {
	Position int `json:"position" validate:"required"`
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

func (r *createPlaylistRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

func (r *editPlaylistRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

func (r *reorderPlaylistRequest) bind(c echo.Context) error {	
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}