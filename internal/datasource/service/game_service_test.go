package service

import (
	"testing"
	"tictactoe/domain/domainmodel"
)

func TestNextMoveWins(t *testing.T) {
	s := &gameService{}
	game := domainmodel.Game{
		Board: domainmodel.Board{
			{domainmodel.Nought, domainmodel.Nought, domainmodel.Empty},
			{domainmodel.Cross, domainmodel.Cross, domainmodel.Empty},
			{domainmodel.Empty, domainmodel.Empty, domainmodel.Cross},
		},
	}

	result, err := s.NextMove(game)
	if err != nil {
		t.Fatal(err)
	}
	if result.Board[0][2] != domainmodel.Nought {
		t.Errorf("expected Nought at [0][2], got %v", result.Board)
	}
}
