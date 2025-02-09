package models

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID    `db:"id"`
	UserID      uuid.UUID    `db:"user_id"`
	PlaylistName        string    `db:"name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Songs	   []Song    `db:"songs"`
}

type PlaylistSong struct {
	ID        uuid.UUID `db:"id"`
	PlaylistID uuid.UUID `db:"playlist_id"`
	SongID    uuid.UUID `db:"song_id"`
	Position  int       `db:"position"`
	CreatedAt time.Time `db:"created_at"`
}