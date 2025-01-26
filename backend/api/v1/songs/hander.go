package songs

import (
	"backend/pkg/songsLib"

	"github.com/labstack/echo/v4"
)

type SongHandler struct {
	songGetter *songsLib.SongGetter
}

func NewHandler() (*SongHandler, error) {
	songGetter, err := songsLib.NewSongGetter()
	if err != nil {
		return nil, err
	}
	return &SongHandler{
		songGetter: songGetter,
	}, nil
}

func (h *SongHandler) Register(group *echo.Group) {
	group.POST("/url", h.GetSongURL)
	group.POST("/details", h.GetSongDetails)
}

