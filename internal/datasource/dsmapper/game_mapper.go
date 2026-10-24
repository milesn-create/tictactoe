package dsmapper

import (
	"errors"
	"tictactoe/datasource/dsmodel"
	"tictactoe/domain/domainmodel"
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
	return domainmodel.Game{ID: game.ID, Board: board}, nil

}
func ToDatasource(game domainmodel.Game) dsmodel.Game {
	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = int(game.Board[i][j])
		}
	}
	return dsmodel.Game{ID: game.ID, Board: board}
}
