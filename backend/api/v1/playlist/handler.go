package playlist

import (
	"os"
	"errors"
	"github.com/labstack/echo/v4"
	"backend/store"
	"backend/pkg/songsLib"
	"github.com/golang-jwt/jwt/v4"
	auth "backend/auth"
	echojwt "github.com/labstack/echo-jwt"
)

type PlaylistHandler struct {
	playlistStore Store
	// TODO: implement interface for others
	songStore *store.SongStore
	songGetter *songsLib.SongGetter
	jwtSecret []byte
}
// TODO: implement interface for others
func NewHandler(playlistStore Store, songStore *store.SongStore, songGetter *songsLib.SongGetter) (*PlaylistHandler, error) {
	secret, ok := os.LookupEnv("Signing_Key")
	if !ok {
		return nil, errors.New("no secret key ")
	}
	return &PlaylistHandler{
		playlistStore: playlistStore,
		songStore: songStore,
		songGetter: songGetter,
		jwtSecret: []byte(secret),
	}, nil
}

func (h *PlaylistHandler) Register(group *echo.Group) {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.jwtSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	})

	playlists := group.Group("/playlists", jwtMiddleware)
	
	playlists.GET("", h.GetPlaylists)
	playlists.GET("/:id", h.GetPlaylist)
	playlists.POST("", h.CreatePlaylist)
	playlists.POST("/:id/song/add", h.AddSongToPlaylist)
	playlists.PUT("/:id", h.EditPlaylist)
	playlists.PUT("/:id/song/reorder", h.ReorderPlaylist)
	playlists.DELETE("/:id", h.RemovePlaylist)
	playlists.DELETE("/:id/song/:song_id", h.RemoveSong)
}