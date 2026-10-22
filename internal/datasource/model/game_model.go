package model

import "github.com/google/uuid"

type Game struct {
	ID    uuid.UUID `db:"id"`
	Board [][]int   `db:"board`
}
