package datasource

import (
	"testing"
	"tictactoe/domain"
)

func TestNextMoveWins(t *testing.T) {
	s := &gameService{}
	game := domain.Game{
		Board: domain.Board{
			{domain.Nought, domain.Nought, domain.Empty},
			{domain.Cross, domain.Cross, domain.Empty},
			{domain.Empty, domain.Empty, domain.Cross},
		},
	}

	result, err := s.NextMove(game)
	if err != nil {
		t.Fatal(err)
	}
	if result.Board[0][2] != domain.Nought {
		t.Errorf("expected Nought at [0][2], got %v", result.Board)
	}
}
