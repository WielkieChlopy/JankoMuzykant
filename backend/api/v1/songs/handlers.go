package songs

import (
	"github.com/labstack/echo/v4"
)

func (h *SongHandler) GetSongDetails(c echo.Context) error {
	req := &songRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(400, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	details, err := h.songGetter.GetSongDetails(req.Url)
	if err != nil {
		return c.JSON(500, map[string]string{
			"error": "Failed to get song details: " + err.Error(),
		})
	}

	return c.JSON(200, details)
}
