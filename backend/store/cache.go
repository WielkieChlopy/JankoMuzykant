package store

import (
	"backend/models"

	"github.com/jmoiron/sqlx"
)


type CacheStore struct {
	db *sqlx.DB
}

func NewCacheStore(db *sqlx.DB) *CacheStore {
	return &CacheStore{
		db: db,
	}
}

func (ss CacheStore) IsSongInCache(songId string, source string) (bool, error) {
	var exists bool
	err := ss.db.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM songs_cache sc
			WHERE sc.song_id = $1 AND sc.source = $2
		)`, songId, source)
	return exists, err
}	

func (ss CacheStore) IsSongExpired(songId string, source string) (bool, error) {
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

func (ss CacheStore) GetSong(songId string, source string) (*models.SongMapping, error) {
	song := &models.SongMapping{}
	err := ss.db.Get(song, `
		SELECT * FROM songs_cache 
		WHERE song_id = $1 AND source = $2`, songId, source)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (ss CacheStore) InsertSong(song *models.SongMapping) error {
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