package songs

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

type SongHandler struct{}

func NewHandler() (*SongHandler, error) {
	return &SongHandler{}, nil
}

func (h *SongHandler) Register(group *echo.Group) {
	group.POST("/url", h.GetSongURL)
}

func (h *SongHandler) GetSongURL(c echo.Context) error {
	req := &songRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(400, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	fmt.Println(req.Url)
	cmd := exec.Command("/app/bin/yt-dlp", "-f", "bestaudio", "-g", req.Url)
	output, err := cmd.Output()
	if err != nil {
		return c.JSON(500, map[string]string{
			"error": "Failed to get audio URL: " + err.Error(),
		})
	}

	audioURL := strings.TrimSpace(string(output))

	return c.JSON(200, map[string]string{
		"url": audioURL,
	})
}
