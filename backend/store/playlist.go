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

//TODO: sprawdz czy potrzebna funkcja
//func (ps PlaylistStore) PlaylistExists(userId uuid.UUID) (bool, error) {
//	var exists bool
//	err := ps.db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM playlist WHERE user_id = $1)", userId)
//	return exists, err
//}

func (ps PlaylistStore) CreatePlaylistForUser(playlist *models.Playlist) (*models.Playlist, error) {
	err := ps.db.QueryRow("INSERT INTO playlists (user_id, name) values ($1, $2) RETURNING id, created_at, updated_at", playlist.UserID, playlist.PlaylistName).Scan(&playlist.ID, &playlist.CreatedAt, &playlist.UpdatedAt)

	if err != nil {
        return nil, err
    }
	
	return playlist, nil
}

func (ps PlaylistStore) GetPlaylistsForUser(userID uuid.UUID) ([]models.Playlist, error) {
	playlists := []models.Playlist{}
	err := ps.db.Select(&playlists, "SELECT * FROM playlists WHERE user_id = $1", userID)	
	return playlists, err
}

func (ps PlaylistStore) GetPlaylistWithSongs(playlistId uuid.UUID, userId uuid.UUID) (*models.Playlist, error) {
	playlist := &models.Playlist{}
	err := ps.db.Get(playlist, "SELECT * FROM playlists WHERE id = $1 AND user_id = $2", playlistId, userId)
	if err != nil {
		//return nil, err
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	songs := []models.Song{}
	//err = ps.db.Select(&songs, `SELECT s.* FROM song s 
	//INNER JOIN playlist_song ps ON s.id = ps.song_id WHERE ps.playlist_id = $1`, playlistId)
	err = ps.db.Select(&songs, "SELECT * FROM song WHERE id IN (SELECT song_id FROM playlist_song WHERE playlist_id = $1)", playlistId)
	if err != nil {
		//return nil, err
		return nil, fmt.Errorf("failed to get songs: %w", err)
	}

	return playlist, nil
}

func (ps PlaylistStore) AddSongToPlaylist(userID uuid.UUID, songID string) error {
	fmt.Println("Adding song to playlist for user", userID)

	tx, err := ps.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	playlist := &models.Playlist{}
	err = tx.Get(playlist, "SELECT id FROM playlists WHERE user_id = $1", userID)
	if err != nil {
		fmt.Println("Error selecting playlistId", err)
		return err
	}

	songId := uuid.UUID{}
	err = tx.Get("SELECT id FROM song WHERE song_id = $1", songID)
	if err != nil {
		fmt.Println("Error selecting songId", err)
		return fmt.Errorf("song with song_id %s not found", songID)
	}

	_, err = tx.Exec("INSERT INTO playlist_song (playlist_id, song_id) VALUES ($1, $2)", playlist.ID, songId)
	if err != nil {
		fmt.Println("Error inserting song into playlist", err)
	}
	return tx.Commit()
}

func (ps PlaylistStore) UpdatePlaylist(playlistName string, userID uuid.UUID, playlistId uuid.UUID) error {
	_, err := ps.db.Exec("UPDATE playlists SET name = $1, updated_at = NOW() WHERE user_id = $2 AND id = $3", playlistName, userID, playlistId)
	if err != nil {
		return err
	}
	return nil
}
// TODO: obsługa zmiany pozycji utworów w playliście
func (ps PlaylistStore) ReorderSongsInPlaylist(playlistId uuid.UUID, songId uuid.UUID, position int) error {
	_, err := ps.db.Exec("UPDATE playlist_song SET position = $1 WHERE playlist_id = $2 AND song_id = $3", position, playlistId, songId)
	return err
}

func (ps PlaylistStore) RemoveSongFromPlaylist(playlistId uuid.UUID, songId uuid.UUID) error {
	_, err := ps.db.Exec("DELETE FROM playlist_song WHERE playlist_id = $1 AND song_id = $2", playlistId, songId)
	return err
}

func (ps PlaylistStore) RemovePlaylist(playlistId uuid.UUID, userId uuid.UUID) error {
	tx, err := ps.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//_, err = tx.Exec("DELETE FROM playlist_song WHERE playlist_id = $1", playlistId)
	//if err != nil {
	//	fmt.Println("Error deleting playlist from playlist_song", err)
	//	return err
	//}
	
	result, err := tx.Exec("DELETE FROM playlists WHERE id = $1 AND user_id = $2", playlistId, userId)
	if err != nil {
		fmt.Println("Error deleting playlist from playlists", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("playlist with id %s not found", playlistId)
	}
		
	return tx.Commit()
}
