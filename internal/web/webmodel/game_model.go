package webmodel

import "github.com/google/uuid"

type Game struct {
	ID    uuid.UUID `json:"id"`
	Board [][]int   `json:"board"`
}
type GameResponse struct {
	ID             uuid.UUID  `json:"id"`
	Board          [][]int    `json:"board"`
	State          string     `json:"state"`
	VsComputer     bool       `json:"vs_computer"`
	TurnPlayerID   *uuid.UUID `json:"turn_player_id,omitempty"`
	WinnerID       *uuid.UUID `json:"winner_id,omitempty"`
	CrossPlayerID  *uuid.UUID `json:"cross_player_id,omitempty"`
	NoughtPlayerID *uuid.UUID `json:"nought_player_id,omitempty"`
}
type CreateGameRequest struct {
	VsComputer bool `json:"vs_computer"`
}
