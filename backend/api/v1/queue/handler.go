package queue

import (
	"github.com/labstack/echo/v4"
)

type QueueHandler struct {
}

func NewHandler() (*QueueHandler, error) {
	return &QueueHandler{}, nil
}

func (h *QueueHandler) Register(group *echo.Group) {
	group.POST("/add", h.AddSong)
	group.POST("/next", h.NextSong)
}
