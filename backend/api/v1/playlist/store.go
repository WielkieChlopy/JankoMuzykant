package playlist

import (
	"backend/models"

	"github.com/google/uuid"
)

type Store interface {
	GetPlaylistsForUser(uuid.UUID) ([]models.Playlist, error)
	GetPlaylistWithSongs(uuid.UUID) (*models.Playlist, error)
	CreatePlaylistForUser(*models.Playlist) (*models.Playlist, error)
	AddSongToPlaylist(uuid.UUID, uuid.UUID) error
	EditPlaylistDetails(*models.Playlist) error
	ReorderSongsInPlaylist(uuid.UUID, uuid.UUID, int) error
	RemovePlaylist(uuid.UUID) error
	RemoveSongFromPlaylist(uuid.UUID, uuid.UUID) error
	//GetByID(uuid.UUID) (*models.Playlist, error)
	//GetByUserID(uuid.UUID) ([]models.Playlist, error)
	//Create(*models.Playlist) (*models.Playlist, error)
	//Update(*models.Playlist) (*models.Playlist, error)
	//Delete(uuid.UUID) error
}