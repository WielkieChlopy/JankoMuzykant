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
	//Songs	   []Song    `db:"-"`
}

type PlaylistSong struct {
	PlaylistID uuid.UUID `db:"playlist_id"`
	SongID    string `db:"song_id"`
	Position  int       `db:"position"`
	CreatedAt time.Time `db:"created_at"`
}