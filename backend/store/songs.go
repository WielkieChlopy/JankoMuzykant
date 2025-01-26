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
