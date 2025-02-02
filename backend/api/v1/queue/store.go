package queue

import (
	"backend/models"

	"github.com/google/uuid"
)

type Store interface {
	QueueExists(uuid.UUID) (bool, error)
	InitQueue(uuid.UUID) error
	AddSongToQueue(uuid.UUID, uuid.UUID) error
	NextSong(uuid.UUID) (*models.Song, error)
	GetQueue(uuid.UUID) ([]models.Song, error)
	ClearQueue(uuid.UUID) error
	RemoveSong(uuid.UUID, uuid.UUID) error
}
