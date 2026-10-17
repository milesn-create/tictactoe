package domain

import "github.com/google/uuid"

type Game struct {
	ID    uuid.UUID
	Board Board
}
