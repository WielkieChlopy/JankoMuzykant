package store

import (
	"backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type QueueStore struct {
	db *sqlx.DB
}

func NewQueueStore(db *sqlx.DB) *QueueStore {
	return &QueueStore{
		db: db,
	}
}

func (ss QueueStore) QueueExists(userId uuid.UUID) (bool, error) {
	var exists bool
	err := ss.db.Get(&exists, "select exists (select 1 from queue where user_id = $1)", userId)
	return exists, err
}

func (ss QueueStore) InitQueue(userId uuid.UUID) error {
	_, err := ss.db.Exec("insert into queue (user_id, next_position, current_position) values ($1, $2, $3) on conflict (user_id) do nothing", userId, 0, 0)
	return err
}

func (ss QueueStore) AddSongToQueue(queueId uuid.UUID, songId uuid.UUID) error {
	tx, err := ss.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var position int
	err = ss.db.Get(&position, "select next_position from queue where id = $1", queueId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into queue_song (queue_id, song_id, position) values ($1, $2, $3)", queueId, songId, position)
	if err != nil {
		return err
	}

	_, err = tx.Exec("update queue set next_position = $1 where id = $2", position + 1, queueId)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (ss QueueStore) NextSong(queueId uuid.UUID) (*models.Song, error) {
	song := &models.Song{}
	err := ss.db.Get(song, "select s.* from songs where song_id = (select song_id from queue_song where queue_id = $1 order by position asc limit 1)", queueId)
	return song, err
}

func (ss QueueStore) GetQueue(queueId uuid.UUID) ([]models.Song, error) {
	songs := []models.Song{}
	err := ss.db.Select(&songs, `
		SELECT s.* FROM song s
		JOIN queue_song qs ON qs.song_id = s.id
		WHERE qs.queue_id = $1
		ORDER BY qs.position ASC`, queueId)
	return songs, err
}

func (ss QueueStore) ClearQueue(queueId uuid.UUID) error {
	tx, err := ss.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("delete from queue_song where queue_id = $1", queueId)
	if err != nil {
		return err
	}
	_, err = tx.Exec("update queue set current_position = 0, next_position = 0, current_song_id = null where user_id = (select user_id from queue where id = $1)", queueId)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (ss QueueStore) RemoveSong(queueId uuid.UUID, songId uuid.UUID) error {
	_, err := ss.db.Exec("delete from queue_song where queue_id = $1 and song_id = $2", queueId, songId)
	return err
}

func (ss QueueStore) GetNextPosition(queueId uuid.UUID) (int, error) {
	var position int
	err := ss.db.Get(&position, "select next_position from queue where id = $1", queueId)
	return position, err
}


