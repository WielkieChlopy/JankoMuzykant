package queue

import (
	"os"
	"backend/utils"
	"net/http"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt"
	auth "backend/auth"
)

type QueueHandler struct {
	store Store
	jwtSecret []byte
}

func NewHandler(store Store) (*QueueHandler, error) {
	secret, ok := os.LookupEnv("Signing_Key")
	if !ok {
		return nil, errors.New("no secret key ")
	}
	return &QueueHandler{
		store: store,
		jwtSecret: []byte(secret),
	}, nil
}

func (h *QueueHandler) Register(group *echo.Group) {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.jwtSecret),
	})

	queue := group.Group("/queue", jwtMiddleware)
	queue.POST("/add", h.AddSong)
	queue.POST("/next", h.NextSong)
	queue.POST("/clear", h.ClearQueue)
	queue.POST("/remove", h.RemoveSong)
	queue.GET("/get", h.GetQueue)
}

func (h *QueueHandler) ensureQueueExists(userID uuid.UUID) error {
	exists, err := h.store.QueueExists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return h.store.InitQueue(userID)
	}
	return nil
}

func (h *QueueHandler) AddSong(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}
	req := &addSongRequest{}

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.store.AddSong(userID, req.SongID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *QueueHandler) NextSong(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.store.NextSong(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *QueueHandler) ClearQueue(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.store.ClearQueue(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *QueueHandler) RemoveSong(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}
	req := &removeSongRequest{}

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.store.RemoveSong(userID, req.SongID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *QueueHandler) GetQueue(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	songs, err := h.store.GetQueue(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, songs)
}
