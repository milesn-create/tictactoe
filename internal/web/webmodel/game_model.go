package webmodel

import "github.com/google/uuid"

type Game struct {
	ID    uuid.UUID `json:"id"`
	Board [][]int   `json:"board"`
}
