package domainmodel

import "github.com/google/uuid"

type LeaderboardEntry struct {
	UserID  uuid.UUID
	Login   string
	WinRate float64
}
