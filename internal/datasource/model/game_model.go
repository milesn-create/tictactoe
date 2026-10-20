package model

import "github.com/google/uuid"

type Game struct {
	ID    uuid.UUID
	Board [][]int
}
