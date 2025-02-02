package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	Id         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	DurationMS int       `db:"duration_ms"`
	URL        string    `db:"url"`
	Source     string    `db:"source"`
	SongID     string    `db:"song_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type SongMapping struct {
	SongID     string    `db:"song_id"`
	Source     string    `db:"source"`
	SongURL    string    `db:"song_url"`
	PlayURL    string    `db:"play_url"`
	DurationMS int       `db:"duration_ms"`
	Title      string    `db:"title"`
	ExpiresAt  time.Time `db:"expires_at"`
	CreatedAt  time.Time `db:"created_at"`
}
