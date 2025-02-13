package playlist

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type addSongRequest struct {
	URL string `json:"url" validate:"required"`
	PlaylistID uuid.UUID `json:"playlist_id" validate:"required"`
}

type removeSongRequest struct {
	SongID uuid.UUID `json:"song_id" validate:"required"`
}

type createPlaylistRequest struct {
	Name string `json:"name" validate:"required"`
}

type editPlaylistRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
	Name string `json:"name"`
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

func (r *removeSongRequest) bind(c echo.Context) error {
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

func (r editPlaylistRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}