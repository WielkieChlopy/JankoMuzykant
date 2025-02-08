package store

import (
	"backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SongStore struct {
	db *sqlx.DB
}

func NewSongStore(db *sqlx.DB) *SongStore {
	return &SongStore{
		db: db,
	}
}

func (ss SongStore) GetSong(songId uuid.UUID) (*models.Song, error) {
	song := &models.Song{}
	err := ss.db.Get(song, "SELECT * FROM song WHERE id = $1", songId)
	return song, err
}

func (ss SongStore) CreateSong(song *models.Song) (*models.Song, error) {
	err := ss.db.Get(song, `
		INSERT INTO song (title, duration_ms, url, source, song_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		song.Title, song.DurationMS, song.URL, song.Source, song.SongID)
	return song, err
}

func (ss SongStore) DeleteSong(songId uuid.UUID) error {
	_, err := ss.db.Exec("DELETE FROM song WHERE id = $1", songId)
	return err
}

func (ss SongStore) UpdateSong(song *models.Song) error {
	_, err := ss.db.Exec(`
		UPDATE song 
		SET title = $1, duration_ms = $2, url = $3, source = $4, song_id = $5, updated_at = NOW()
		WHERE id = $6`,
		song.Title, song.DurationMS, song.URL, song.Source, song.SongID, song.Id)
	return err
}

func (ss SongStore) ListSongs() ([]models.Song, error) {
	songs := []models.Song{}
	err := ss.db.Select(&songs, "SELECT * FROM song ORDER BY title")
	return songs, err
}

func (ss SongStore) GetSongBySourceAndSongId(source string, songId string) (*models.Song, error) {
	song := &models.Song{}
	err := ss.db.Get(song, "SELECT * FROM song WHERE source = $1 AND song_id = $2", source, songId)
	return song, err
}
