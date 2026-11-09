package dsmodel

import "github.com/google/uuid"

type LeaderboardRow struct {
	UserID  uuid.UUID `db:"user_id"`
	Login   string    `db:"login"`
	WinRate float64   `db:"win_rate"`
}
