package playlist

import (
	auth "backend/auth"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *PlaylistHandler) GetPlaylists(c echo.Context) error {
	fmt.Println("Getting playlists")
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	playlists, err := h.playlistStore.GetPlaylistsForUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, playlists)
}

func (h *PlaylistHandler) GetPlaylist(c echo.Context) error {
	fmt.Println("Getting playlist")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	playlist, err := h.playlistStore.GetPlaylistWithSongs(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, playlist)
}

func (h *PlaylistHandler) CreatePlaylist(c echo.Context) error {
	fmt.Println("Creating playlist")
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	req := &createPlaylistRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	playlist := &models.Playlist {
		PlaylistName:   req.Name,
		UserID: userID,
	}
	playlist, err = h.playlistStore.CreatePlaylistForUser(playlist)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, playlist)
}

func (h *PlaylistHandler) AddSongToPlaylist(c echo.Context) error {
    fmt.Println("Adding song to playlist")
	//TODO: nie wiem czy jesli UserID jest nieuzywany to moze byc bez indexu czy calkiem elo
    _, err := auth.UserIDFromToken(c)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, utils.NewError(err))
    }

    req := &addSongRequest{}
    if err := req.bind(c); err != nil {
        return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }

    u, err := url.Parse(req.URL)
    if err != nil {
        return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }

    source := h.songGetter.GetSongSource(u)
    sourceSongID := h.songGetter.GetSongId(u, source)
    songID, err := uuid.Parse(sourceSongID)
    if err != nil {
        return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
    }

    if err := h.playlistStore.AddSongToPlaylist(req.PlaylistID, songID); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    return c.JSON(http.StatusOK, "Song added to playlist")
}

func (h *PlaylistHandler) EditPlaylist(c echo.Context) error {
	fmt.Println("Editing playlist")
	_, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	req := &editPlaylistRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	playlist := &models.Playlist {
		PlaylistName: req.Name,
		ID: req.ID,
	}

	if err := h.playlistStore.EditPlaylistDetails(playlist); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Playlist edited")
}
// TODO:
func (h *PlaylistHandler) ReorderPlaylist(c echo.Context) error {
	return nil
}

func (h *PlaylistHandler) RemovePlaylist(c echo.Context) error {
	fmt.Println("Removing playlist")
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	if err := h.playlistStore.RemovePlaylist(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Playlist removed")
}

func (h *PlaylistHandler) RemoveSong(c echo.Context) error {
	fmt.Println("Removing song from playlist")
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}
	req := &removeSongRequest{}

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.playlistStore.RemoveSongFromPlaylist(userID, req.SongID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Song removed from playlist")
}