package dsmodel

import (
	"github.com/google/uuid"
)

type Game struct {
	ID             uuid.UUID  `db:"id"`
	Board          [][]int    `db:"board"`
	State          int        `db:"state"`
	TurnPlayerID   *uuid.UUID `db:"turn_player_id"`
	WinnerID       *uuid.UUID `db:"winner_id"`
	VsComputer     bool       `db:"vs_computer"`
	CrossPlayerID  *uuid.UUID `db:"cross_player_id"`
	NoughtPlayerID *uuid.UUID `db:"nought_player_id"`
}
