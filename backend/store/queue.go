package store

import (
	"backend/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"errors"
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
	if err := ss.db.Get(&exists, "select exists (select 1 from queue where user_id = $1)", userId); err != nil {
		return false, fmt.Errorf("error checking if queue exists: %w", err)
	}
	return exists, nil
}

func (ss QueueStore) InitQueue(userId uuid.UUID) error {
	_, err := ss.db.Exec("insert into queue (user_id, next_position, current_position) values ($1, $2, $3) on conflict (user_id) do nothing", userId, 0, 0)
	if err != nil {
		return fmt.Errorf("error initializing queue: %w", err)
	}
	return nil
}s

func (ss QueueStore) ChangeCurrentSong(queueId uuid.UUID, songId uuid.UUID) error {
	tx, err := ss.db.Beginx()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	var position int
	err = ss.db.Get(&position, "select current_position from queue where user_id = $1", queueId)
	if err != nil {
		return fmt.Errorf("error getting current position: %w", err)
	}

	_, err = tx.Exec("update queue_song set song_id = $1 where queue_id = $2 and position = $3", songId, queueId, position)
	if err != nil {
		return fmt.Errorf("error updating queue song: %w", err)
	}

	_, err = tx.Exec("update queue set current_song_id = $1 where user_id = $2", songId, queueId)
	if err != nil {
		return fmt.Errorf("error updating queue: %w", err)
	}

	return tx.Commit()
}

func (ss QueueStore) AddSongToQueue(queueId uuid.UUID, songId uuid.UUID) error {
	fmt.Println("Adding song to queue", queueId, songId)
	tx, err := ss.db.Beginx()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	var position int
	err = tx.Get(&position, "select next_position from queue where user_id = $1", queueId)
	if err != nil {
		fmt.Println("Error getting next position", err)
		return fmt.Errorf("error getting next position: %w", err)
	}

	_, err = tx.Exec("insert into queue_song (queue_id, song_id, position) values ($1, $2, $3)", queueId, songId, position)
	if err != nil {
		fmt.Println("Error inserting song into queue", err)
		return fmt.Errorf("error inserting song into queue: %w", err)
	}

	_, err = tx.Exec("update queue set next_position = $1 where user_id = $2", position + 1, queueId)
	if err != nil {
		fmt.Println("Error updating next position", err)
		return fmt.Errorf("error updating next position: %w", err)
	}
	return tx.Commit()
}

func (ss QueueStore) GetNextSong(queueId uuid.UUID) (*models.Song, error) {
	song := &models.Song{}
	err := ss.db.Get(song, "select s.* from song s where s.song_id = (select qs.song_id from queue_song qs where qs.queue_id = $1 order by qs.position asc offset 1 limit 1)", queueId)
	return song, fmt.Errorf("error getting next song: %w", err)
}

func (ss QueueStore) PlayNextSong(queueId uuid.UUID) (*models.Song, error) {
	tx, err := ss.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()
	
	type NextSongPos struct {
		SongId uuid.UUID `db:"song_id"`
		Position int `db:"position"`
	}

	var nextSongPos NextSongPos
	err = tx.Get(&nextSongPos, "select song_id, position from queue_song where queue_id = $1 order by position asc offset 1 limit 1", queueId)
	if err != nil {
		return nil, fmt.Errorf("error getting next song position: %w", err)
	}

	song := &models.Song{}
	err = tx.Get(song, "select * from song where song_id = $1", nextSongPos.SongId)
	if err != nil {
		return nil, fmt.Errorf("error getting next song: %w", err)
	}

	_, err = tx.Exec("update queue set current_song_id = $1, current_position = $2 where user_id = $3", song.Id, nextSongPos.Position, queueId)
	if err != nil {
		return nil, fmt.Errorf("error updating queue: %w", err)
	}

	return song, tx.Commit()
}

func (ss QueueStore) GetQueue(queueId uuid.UUID) ([]models.Song, error) {
	songs := []models.Song{}
	err := ss.db.Select(&songs, `
		SELECT s.* FROM song s
		JOIN queue_song qs ON qs.song_id = s.id
		WHERE qs.queue_id = $1
		ORDER BY qs.position ASC`, queueId)
	return songs, fmt.Errorf("error getting queue: %w", err)
}

func (ss QueueStore) ClearQueue(queueId uuid.UUID) error {
	tx, err := ss.db.Beginx()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec("delete from queue_song where queue_id = $1", queueId)
	if err != nil {
		return fmt.Errorf("error deleting queue song: %w", err)
	}
	_, err = tx.Exec("update queue set current_position = 0, next_position = 0, current_song_id = null where user_id = $1", queueId)
	if err != nil {
		return fmt.Errorf("error updating queue: %w", err)
	}
	return tx.Commit()
}

func (ss QueueStore) RemoveSong(queueId uuid.UUID, songId uuid.UUID) error {
	_, err := ss.db.Exec("delete from queue_song where queue_id = $1 and song_id = $2", queueId, songId)
	return fmt.Errorf("error removing song from queue: %w", err)
}

func (ss QueueStore) GetNextPosition(queueId uuid.UUID) (int, error) {
	var position int
	err := ss.db.Get(&position, "select next_position from queue where user_id = $1", queueId)
	return position, fmt.Errorf("error getting next position: %w", err)
}


