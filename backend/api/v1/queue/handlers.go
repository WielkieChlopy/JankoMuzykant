package queue

import (
	"database/sql"
	"net/http"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"backend/utils"
	"net/url"
	auth "backend/auth"
	"errors"
	"backend/models"
	"backend/pkg/songsLib"
)

func (h *QueueHandler) ensureQueueExists(userID uuid.UUID) error {
	exists, err := h.queueStore.QueueExists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return h.queueStore.InitQueue(userID)
	}
	return nil
}

func (h *QueueHandler) PlaySong(c echo.Context) error {
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

	u, err := url.Parse(req.URL)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	source := h.songGetter.GetSongSource(u)
	source_songId := h.songGetter.GetSongId(u, source)

	song, details, err := EnsureSongExistance(h, source, source_songId, req, c)
	if err != nil {
		return err
	}

	if err := h.queueStore.ChangeCurrentSong(userID, song.Id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, details)
}

func (h *QueueHandler) AddSong(c echo.Context) error {
	fmt.Println("Adding song")
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

	u, err := url.Parse(req.URL)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	fmt.Println("Parsed URL")

	source := h.songGetter.GetSongSource(u)
	source_songId := h.songGetter.GetSongId(u, source)
	fmt.Println("Got source and songId")
	song, details, err := EnsureSongExistance(h, source, source_songId, req, c)
	if err != nil {
		return err
	}
	fmt.Println("Ensured song existence", err)

	if err := h.queueStore.AddSongToQueue(userID, song.Id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, details)
}

func EnsureSongExistance(h *QueueHandler, source string, source_songId string, req *addSongRequest, c echo.Context) (*models.Song, *songsLib.SongDetails, error) {
	song, err := h.songStore.GetSongBySourceAndSongId(source, source_songId)
	if err == nil {
		fmt.Println("Song already exists")
		return song, nil, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Song does not exist")
		details, err := h.songGetter.GetSongDetails(req.URL)
		if err != nil {
			return nil, nil, echo.NewHTTPError(http.StatusUnprocessableEntity, utils.NewError(err))
		}

		song = &models.Song{
			Title:      details.Title,
			DurationMS: int(details.DurationMS),
			URL:        req.URL,
			Source:     source,
			SongID:     source_songId,
		}
		song,err = h.songStore.CreateSong(song)
		if err != nil {
			return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, utils.NewError(err))
		}

		return song, &details, nil
	} 
	fmt.Println("Error", err)
	return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, utils.NewError(err))
}

func (h *QueueHandler) PlayNextSong(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	song, err := h.queueStore.PlayNextSong(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	details, err := h.songGetter.GetSongDetails(song.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, details)
}


func (h *QueueHandler) GetNextSong(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	song, err := h.queueStore.GetNextSong(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	details, err := h.songGetter.GetSongDetails(song.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, details)
}

func (h *QueueHandler) ClearQueue(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.queueStore.ClearQueue(userID); err != nil {
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

	if err := h.queueStore.RemoveSong(userID, req.SongID); err != nil {
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

	songs, err := h.queueStore.GetQueue(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, songs)
}
