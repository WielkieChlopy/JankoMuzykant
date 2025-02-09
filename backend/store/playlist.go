package store

import (
	"backend/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PlaylistStore struct {
	db *sqlx.DB
}

func NewPlaylistStore(db *sqlx.DB) *PlaylistStore {
	return &PlaylistStore{
		db: db,
	}
}

// stworzenie playlisty
// listownaie playlist
// pobranie playlisty z utworami - pobiera dane playlisty i utwory w niej zawarte
// dodawanie utworów do playlisty
// usuwanie utworów z playlisty
// usuwanie playlisty



// second plan:
// granie utworów z playlisty - GetNextSong
// przelaczanie utworow - kolejny utwor
// przelaczanie utworow - poprzedni utwor
// przelaczanie utworow - losowy utwor
// zmiana kolejnosci utworow w playliscie

func (ps PlaylistStore) CreatePlaylist(playlist *models.Playlist) (*models.Playlist, error) {
	err := ps.db.Get(playlist, `
		INSERT INTO playlist (user_id, name) values ($1, $2) RETURNING *`, playlist.UserID, playlist.PlaylistName)
	return playlist, err
}

func (ps PlaylistStore) ListPlaylistsForUser(userID uuid.UUID) ([]models.Playlist, error) {
	playlists := []models.Playlist{}
	err := ps.db.Select(&playlists, "SELECT * FROM playlist WHERE user_id = $1", userID)	
	return playlists, err
}

func (ps PlaylistStore) GetPlaylistWithSongs(playlistId uuid.UUID) (*models.Playlist, error) {
	playlist := &models.Playlist{}
	err := ps.db.Get(playlist, "SELECT * FROM playlist WHERE id = $1", playlistId)
	if err != nil {
		return nil, err
	}

	songs := []models.Song{}
	err = ps.db.Select(&songs, "SELECT * FROM song WHERE id IN (SELECT song_id FROM playlist_song WHERE playlist_id = $1)", playlistId)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (ps PlaylistStore) AddSongToPlaylist(playlistId uuid.UUID, songId uuid.UUID) error {
	fmt.Println("Addin song to playlist", playlistId, songId)

	_, err := ps.db.Exec("INSERT INTO playlist_song (playlist_id, song_id) values ($1, $2)", playlistId, songId)
	if err != nil {
		return err
	}

	return nil
}

func (ps PlaylistStore) RemoveSongFromPlaylist(playlistId uuid.UUID, songId uuid.UUID) error {
	_, err := ps.db.Exec("DELETE FROM playlist_song WHERE playlist_id = $1 AND song_id = $2", playlistId, songId)
	return err
}

func (ps PlaylistStore) RemovePlaylist(playlistId uuid.UUID) error {
	_, err := ps.db.Exec("DELETE FROM playlist WHERE id = $1", playlistId)
	return err
}
