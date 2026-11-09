package webmodel

import "github.com/google/uuid"

type LeaderboardResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Login   string    `json:"login"`
	WinRate float64   `json:"win_rate"`
}
