package playlist

import (
	auth "backend/auth"
	"backend/models"
	"backend/pkg/songsLib"
	"backend/utils"
	"database/sql"
	"errors"
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
	_, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	playlist, err := h.playlistStore.GetPlaylistWithSongs(playlistID)
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

func (h *PlaylistHandler) AddSong(c echo.Context) error {
    fmt.Println("Adding song to playlist")
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
	fmt.Println("Parsed URL")

	source := h.songGetter.GetSongSource(u)
	sourceSongID := h.songGetter.GetSongId(u, source)
	fmt.Println("Got source and songId")

	song, details, err := EnsureSongExistance(h, source, sourceSongID, req, c)
	if err != nil {
		return err
	}
	fmt.Println("Ensured song existence", err)

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

    if err := h.playlistStore.AddSongToPlaylist(playlistID, song.Id); err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewError(err))
    }

	fmt.Println("Added song to playlist")
    return c.JSON(http.StatusOK, details)
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

func (h *PlaylistHandler) ReorderPlaylist(c echo.Context) error {
	fmt.Println("Reordering playlist")
	_, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	
	req := &reorderPlaylistRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	songID, err := uuid.Parse(c.Param("song_id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.playlistStore.ReorderSongsInPlaylist(playlistID, songID, req.Position); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Playlist reordered")
}

func (h *PlaylistHandler) RemoveSong(c echo.Context) error {
	fmt.Println("Removing song from playlist")
	_, err := auth.UserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewError(err))
	}

	songID, err := uuid.Parse(c.Param("song_id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	playlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := h.playlistStore.RemoveSongFromPlaylist(playlistID, songID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, "Song removed from playlist")
}

func EnsureSongExistance(h *PlaylistHandler, source string, sourceSongID string, req *addSongRequest, c echo.Context) (*models.Song, *songsLib.SongDetails, error) {
	song, err := h.songStore.GetSongBySourceAndSongId(source, sourceSongID)
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
			SongID:     sourceSongID,
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