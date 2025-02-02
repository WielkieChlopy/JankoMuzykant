package queue

import (
	"os"
	"errors"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt"
	"backend/pkg/songsLib"
	"backend/store"
	"github.com/golang-jwt/jwt/v4"
	auth "backend/auth"
)

type QueueHandler struct {
	queueStore *store.QueueStore
	songStore *store.SongStore
	songGetter *songsLib.SongGetter
	jwtSecret []byte
}

func NewHandler(queueStore *store.QueueStore, songStore *store.SongStore, songGetter *songsLib.SongGetter) (*QueueHandler, error) {
	secret, ok := os.LookupEnv("Signing_Key")
	if !ok {
		return nil, errors.New("no secret key ")
	}
	return &QueueHandler{
		queueStore: queueStore,
		songStore: songStore,
		songGetter: songGetter,
		jwtSecret: []byte(secret),
	}, nil
}

func (h *QueueHandler) Register(group *echo.Group) {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.jwtSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	})

	queue := group.Group("/queue", jwtMiddleware)
	queue.POST("", h.AddSong)
	queue.POST("/next", h.NextSong)
	queue.DELETE("", h.ClearQueue)
	queue.POST("/remove", h.RemoveSong)
	queue.GET("", h.GetQueue)
}

