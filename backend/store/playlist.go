package store

import (
	"backend/models"
	"database/sql"
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
	err := ps.db.Select(&playlists, "SELECT * FROM playlists WHERE user_id = $1 ORDER BY updated_at DESC", userID)	
	return playlists, err
}

func (ps PlaylistStore) GetPlaylistWithSongs(playlistId uuid.UUID) ([]models.Song, error) {
	songs := []models.Song{}
	err := ps.db.Select(&songs, `
		SELECT s.* FROM song s 
		JOIN playlist_song ps ON s.id = ps.song_id
		WHERE ps.playlist_id = $1
		ORDER BY ps.position ASC
		`, playlistId)
	return songs, err
}

func (ps PlaylistStore) AddSongToPlaylist(playlistId uuid.UUID, songId uuid.UUID) error {
	fmt.Printf("Adding song to playlist %s", playlistId)
	tx, err := ps.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	var max_position int
	err = tx.Get(&max_position, "SELECT COALESCE(MAX(position)) FROM playlist_song WHERE playlist_id = $1", playlistId)
	if err != nil {
		fmt.Println("Error selecting position", err)
		return err
	}

	new_position := max_position + 1
	_, err = tx.Exec("INSERT INTO playlist_song (playlist_id, song_id, position) VALUES ($1, $2, $3)", playlistId, songId, new_position)
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

func (ps PlaylistStore) ReorderSongsInPlaylist(playlistId uuid.UUID, songId uuid.UUID, newPosition int) error {
	tx, err := ps.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	// blokada playlisty to update
	if _, err := tx.Exec("SELECT 1 FROM playlists WHERE id = $1 FOR UPDATE", playlistId); err != nil {
		return fmt.Errorf("failed to lock playlist with ID=%v: %w", playlistId, err)
	}

	var currentPosition int
	err = tx.Get(&currentPosition, "SELECT position FROM playlist_song WHERE playlist_id = $1 AND song_id = $2", playlistId, songId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("song with id %s not found in playlist with id %s", songId, playlistId)
	} else if err != nil {
		return fmt.Errorf("failed to get current position: %w", err)
	}

	if currentPosition == newPosition {
		fmt.Println("Current position is equal to new position")
		return nil
	}

	var maxPosition int
	if err := tx.Get(&maxPosition, "SELECT COALESCE(MAX(position), 0) FROM playlist_song WHERE playlist_id = $1", playlistId); err != nil {
		return fmt.Errorf("failed to get max position: %w", err)
	}

	if newPosition < 0 || newPosition > maxPosition {
		return fmt.Errorf("new position %d is out of range [0, %d]", newPosition, maxPosition)
	}

	if newPosition > currentPosition {
		_, err = tx.Exec(`UPDATE playlist_song SET position = position - 1 WHERE playlist_id = $1 AND position > $2 AND position <= $3`, playlistId, currentPosition, newPosition)
	} else {
		_, err = tx.Exec(`UPDATE playlist_song SET position = position + 1 WHERE playlist_id = $1 AND position >= $2 AND position < $3`, playlistId, newPosition, currentPosition)
	}

	if err != nil {
		return fmt.Errorf("failed to update positions: %w", err)
	}

	if _, err := tx.Exec("UPDATE playlist_song SET position = $1 WHERE playlist_id = $2 AND song_id = $3", newPosition, playlistId, songId); err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
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
