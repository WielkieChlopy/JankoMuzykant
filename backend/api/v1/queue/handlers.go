package queue

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"backend/utils"
	"net/url"
	auth "backend/auth"
	"errors"
	"backend/models"
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

	u, err := url.Parse(req.URL)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	source := h.songGetter.GetSongSource(u)
	source_songId := h.songGetter.GetSongId(u, source)

	song, err := EnsureSongExistance(h, source, source_songId, req, c)
	if err != nil {
		return err
	}

	if err := h.queueStore.AddSongToQueue(userID, song.Id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

func EnsureSongExistance(h *QueueHandler, source string, source_songId string, req *addSongRequest, c echo.Context) (*models.Song, error) {
	song, err := h.songStore.GetSongBySourceAndSongId(source, source_songId)
	if err == nil {
		return song, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		details, err := h.songGetter.GetSongDetails(req.URL)
		if err != nil {
			return nil, c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
		}

		song = &models.Song{
			Title:      details.Title,
			DurationMS: int(details.DurationMS),
			URL:        req.URL,
			Source:     source,
			SongID:     source_songId,
		}
		err = h.songStore.CreateSong(song)
		if err != nil {
			return nil, c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
	} 
	return nil, c.JSON(http.StatusInternalServerError, utils.NewError(err))
}

func (h *QueueHandler) NextSong(c echo.Context) error {
	userID, err := auth.UserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.ensureQueueExists(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	song, err := h.queueStore.NextSong(userID)
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
