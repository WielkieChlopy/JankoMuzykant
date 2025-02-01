package v1

import (
	"backend/api/v1/songs"
	"backend/api/v1/user"
	"backend/store"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserHandler user.UserHandler
	SongHandler songs.SongHandler
}

func NewHandler(userS *store.UserStore, songS *store.SongStore) (*Handler, error) {
	uh, err := user.NewHandler(userS)
	if err != nil {
		return nil, err
	}

	sh, err := songs.NewHandler(songS)
	if err != nil {
		return nil, err
	}

	return &Handler{
		UserHandler: *uh,
		SongHandler: *sh,
	}, nil
}

func (h *Handler) Register(group *echo.Group) {
	h.UserHandler.Register(group)
	h.SongHandler.Register(group)
}
