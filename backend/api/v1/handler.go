package v1

import (
	"backend/api/v1/user"
	"backend/store"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserHandler user.UserHandler
}

func NewHandler(userS *store.UserStore) (*Handler, error) {
	uh, err := user.NewHandler(userS)
	if err != nil {
		return nil, err
	}

	return &Handler{
		UserHandler: *uh,
	}, nil
}

func (h *Handler) Register(group *echo.Group) {
	h.UserHandler.Register(group)
}
