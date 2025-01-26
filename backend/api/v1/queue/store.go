package queue

import (
	"backend/models"

	"github.com/google/uuid"
)

type Store interface {
	AddSong(uuid.UUID, uuid.UUID) error
	NextSong(uuid.UUID) error
	GetUserQueue(uuid.UUID) ([]models.Song, error)
}
