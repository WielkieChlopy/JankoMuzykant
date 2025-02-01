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

func (ss SongStore) InitQueue(userId uuid.UUID) error {
	_, err := ss.db.Exec("insert into queue (user_id) values ($1) on conflict (user_id) do nothing", userId)
	return err
}

func (ss SongStore) AddSong(queueId uuid.UUID, songId uuid.UUID) error {
	_, err := ss.db.Exec("insert into queue_song (queue_id, song_id) values ($1, $2)", queueId, songId)
	return err
}

func (ss SongStore) NextSong(queueId uuid.UUID) error {
	// Get the first song in queue and delete it
	_, err := ss.db.Exec("delete from queue_song where queue_id = $1 and id = (select id from queue_song where queue_id = $1 limit 1)", queueId)
	return err
}

func (ss SongStore) GetQueue(queueId uuid.UUID) ([]models.Song, error) {
	songs := []models.Song{}
	err := ss.db.Select(&songs, `
		SELECT s.* FROM song s
		JOIN queue_song qs ON qs.song_id = s.id
		WHERE qs.queue_id = $1
		ORDER BY qs.id ASC`, queueId)
	return songs, err
}

func (ss SongStore) ClearQueue(queueId uuid.UUID) error {
	_, err := ss.db.Exec("delete from queue_song where queue_id = $1", queueId)
	return err
}

func (ss SongStore) RemoveSong(queueId uuid.UUID, songId uuid.UUID) error {
	_, err := ss.db.Exec("delete from queue_song where queue_id = $1 and song_id = $2", queueId, songId)
	return err
}

func (ss SongStore) IsSongInCache(songId string, source string) (bool, error) {
	var exists bool
	err := ss.db.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM songs_cache sc
			WHERE sc.song_id = $1 AND sc.source = $2
		)`, songId, source)
	return exists, err
}	

func (ss SongStore) IsSongExpired(songId string, source string) (bool, error) {
	var expired bool
	err := ss.db.Get(&expired, `
		SELECT CASE 
			WHEN expires_at < NOW() THEN true 
			ELSE false 
		END 
		FROM songs_cache 
		WHERE song_id = $1 AND source = $2`, songId, source)
	return expired, err
}

func (ss SongStore) GetSong(songId string, source string) (*models.SongMapping, error) {
	song := &models.SongMapping{}
	err := ss.db.Get(song, `
		SELECT * FROM songs_cache 
		WHERE song_id = $1 AND source = $2`, songId, source)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (ss SongStore) InsertSong(song *models.SongMapping) error {
	_, err := ss.db.Exec(`
		INSERT INTO songs_cache (song_id, source, song_url, play_url, duration_ms, title, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (song_id, source) 
		DO UPDATE SET 
			song_url = EXCLUDED.song_url,
			play_url = EXCLUDED.play_url,
			duration_ms = EXCLUDED.duration_ms,
			title = EXCLUDED.title,
			expires_at = EXCLUDED.expires_at`,
		song.SongID, song.Source, song.SongURL, song.PlayURL, song.DurationMS, song.Title, song.ExpiresAt)
	return err
}
