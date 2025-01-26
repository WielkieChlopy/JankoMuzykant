package queue

import (
	"backend/models"

	"github.com/google/uuid"
)

type Store interface {
	QueueExists(uuid.UUID) (bool, error)
	InitQueue(uuid.UUID) error
	AddSong(uuid.UUID, uuid.UUID) error
	NextSong(uuid.UUID) error
	GetQueue(uuid.UUID) ([]models.Song, error)
	ClearQueue(uuid.UUID) error
	RemoveSong(uuid.UUID, uuid.UUID) error
}
