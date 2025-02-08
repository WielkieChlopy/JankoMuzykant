package v1

import (
	"backend/api/v1/songs"
	"backend/api/v1/user"
	"backend/store"
	"backend/pkg/songsLib"
	"backend/api/v1/queue"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserHandler user.UserHandler
	SongHandler songs.SongHandler
	QueueHandler queue.QueueHandler
}

func NewHandler(userS *store.UserStore, songS *store.SongStore, queueS *store.QueueStore, cacheS *store.CacheStore) (*Handler, error) {
	uh, err := user.NewHandler(userS)
	if err != nil {
		return nil, err
	}

	songGetter, err := songsLib.NewSongGetter(queueS, cacheS)
	if err != nil {
		return nil, err
	}

	sh, err := songs.NewHandler(songGetter)
	if err != nil {
		return nil, err
	}

	qh, err := queue.NewHandler(queueS, songS, songGetter)
	if err != nil {
		return nil, err
	}

	return &Handler{
		UserHandler: *uh,
		SongHandler: *sh,
		QueueHandler: *qh,
	}, nil
}

func (h *Handler) Register(group *echo.Group) {
	h.UserHandler.Register(group)
	h.SongHandler.Register(group)
	h.QueueHandler.Register(group)
}
