package models

import "github.com/google/uuid"

type Song struct {
	Id         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	DurationMS int       `db:"duration_ms"`
	PlayUrl    string    `db:"play_url"`
	URL        string    `db:"url"`
	UserId     uuid.UUID `db:"user_id"`
}
