package user

import (
	"backend/models"

	"github.com/google/uuid"
)

type Store interface {
	GetByID(uuid.UUID) (*models.User, error)
	GetByUsername(string) (*models.User, error)
	Create(*models.User) (*models.User, error)
	Update(*models.User) (*models.User, error)
	Delete(uuid.UUID) error
}
