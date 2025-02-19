package playlist

import (
	auth "backend/auth"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"net/url"
	//"strings"

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
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	playlist, err := h.playlistStore.GetPlaylistWithSongs(id, userID)
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
		UserID:         userID,
		PlaylistName:   req.Name,
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
    userID, err := auth.UserIDFromToken(c)
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

    if err := h.playlistStore.AddSongToPlaylist(userID, sourceSongID); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

    return c.JSON(http.StatusOK, "Song added to playlist")
}

func (h *PlaylistHandler) EditPlaylist(c echo.Context) error {
	fmt.Println("Editing playlist")
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	req := &editPlaylistRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.playlistStore.UpdatePlaylist(req.Name, userID, playlistID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Playlist edited")
}
// TODO:
func (h *PlaylistHandler) ReorderPlaylist(c echo.Context) error {
	return nil
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

func (h *PlaylistHandler) RemovePlaylist(c echo.Context) error {
	fmt.Println("Removing playlist")
	userID, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.playlistStore.RemovePlaylist(playlistID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Playlist removed")
}