package playlist

import (
	"backend/models"

	"github.com/google/uuid"
)

type Store interface {
	GetPlaylistsForUser(uuid.UUID) ([]models.Playlist, error)
	GetPlaylistWithSongs(uuid.UUID, uuid.UUID) (*models.Playlist, error)
	CreatePlaylistForUser(*models.Playlist) (*models.Playlist, error)
	AddSongToPlaylist(uuid.UUID, string) error
	UpdatePlaylist(string, uuid.UUID, uuid.UUID) error
	ReorderSongsInPlaylist(uuid.UUID, uuid.UUID, int) error
	RemovePlaylist(uuid.UUID, uuid.UUID) error
	RemoveSongFromPlaylist(uuid.UUID, uuid.UUID) error
	//GetByID(uuid.UUID) (*models.Playlist, error)
	//GetByUserID(uuid.UUID) ([]models.Playlist, error)
	//Create(*models.Playlist) (*models.Playlist, error)
	//Update(*models.Playlist) (*models.Playlist, error)
	//Delete(uuid.UUID) error
}