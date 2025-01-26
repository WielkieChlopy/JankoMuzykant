package songs

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

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
