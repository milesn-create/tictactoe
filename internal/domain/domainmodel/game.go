package domainmodel

import "github.com/google/uuid"

type Game struct {
	ID           uuid.UUID
	Board        Board
	State        GameState
	TurnPlayerID uuid.UUID
	WinnerID     uuid.UUID
}
