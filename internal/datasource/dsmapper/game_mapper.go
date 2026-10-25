package dsmapper

import (
	"errors"
	"tictactoe/datasource/dsmodel"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
)

var ErrInvalidBoardSize = errors.New("board must be 3x3")

func ToDomain(game dsmodel.Game) (domainmodel.Game, error) {

	if len(game.Board) != 3 {
		return domainmodel.Game{}, ErrInvalidBoardSize

	}
	var board domainmodel.Board
	for i := 0; i < 3; i++ {
		if len(game.Board[i]) != 3 {
			return domainmodel.Game{}, ErrInvalidBoardSize

		}
		for j := 0; j < 3; j++ {
			board[i][j] = domainmodel.Cell(game.Board[i][j])

		}
	}
	state := domainmodel.GameState(game.State)
	return domainmodel.Game{ID: game.ID, Board: board, State: state, TurnPlayerID: ptrToUUID(game.TurnPlayerID), WinnerID: ptrToUUID(game.WinnerID), VsComputer: game.VsComputer, CrossPlayerID: ptrToUUID(game.CrossPlayerID), NoughtPlayerID: ptrToUUID(game.NoughtPlayerID)}, nil

}
func ToDatasource(game domainmodel.Game) dsmodel.Game {
	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = int(game.Board[i][j])
		}
	}
	state := int(game.State)
	return dsmodel.Game{ID: game.ID, Board: board, State: state, TurnPlayerID: uuidToPtr(game.TurnPlayerID), WinnerID: uuidToPtr(game.WinnerID), VsComputer: game.VsComputer, CrossPlayerID: uuidToPtr(game.CrossPlayerID), NoughtPlayerID: uuidToPtr(game.NoughtPlayerID)}
}
func uuidToPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}
func ptrToUUID(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}
