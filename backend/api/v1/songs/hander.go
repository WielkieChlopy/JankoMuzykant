package songs

import (
	"errors"
	"os"

	"backend/pkg/songsLib"

	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt"
)

type SongHandler struct {
	songGetter *songsLib.SongGetter
	jwtSecret []byte
}

func NewHandler() (*SongHandler, error) {
	songGetter, err := songsLib.NewSongGetter()
	if err != nil {
		return nil, err
	}
	secret, ok := os.LookupEnv("Signing_Key")
	if !ok {
		return nil, errors.New("no secret key ")
	}
	return &SongHandler{
		songGetter: songGetter,
		jwtSecret: []byte(secret),
	}, nil
}

func (h *SongHandler) Register(group *echo.Group) {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.jwtSecret),
	})

	songs := group.Group("/songs", jwtMiddleware)

	songs.POST("/details", h.GetSongDetails)
}

