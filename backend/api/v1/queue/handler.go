package queue

import (
	"backend/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QueueHandler struct {
	store Store
}

func NewHandler(store Store) (*QueueHandler, error) {
	return &QueueHandler{
		store: store,
	}, nil
}

func (h *QueueHandler) Register(group *echo.Group) {
	group.POST("/add", h.AddSong)
	group.POST("/next", h.NextSong)
	group.POST("/clear", h.ClearQueue)
	group.POST("/remove", h.RemoveSong)
	group.GET("/get", h.GetQueue)
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
	userID := c.Get("user_id").(uuid.UUID)
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
	userID := c.Get("user_id").(uuid.UUID)

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.store.NextSong(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *QueueHandler) ClearQueue(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.store.ClearQueue(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func (h *QueueHandler) RemoveSong(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
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
	userID := c.Get("user_id").(uuid.UUID)

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	songs, err := h.store.GetQueue(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, songs)
}
