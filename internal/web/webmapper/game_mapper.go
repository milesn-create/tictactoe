package webmapper

import (
	"errors"
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"
)

var (
	ErrInvalidBoardSize = errors.New("board must be 3x3")
	ErrInvalidCellValue = errors.New("cell must be -1, 0 or 1")
)

func ToDomain(game webmodel.Game) (domainmodel.Game, error) {
	if len(game.Board) != 3 {
		return domainmodel.Game{}, ErrInvalidBoardSize
	}
	var board domainmodel.Board
	for i := 0; i < 3; i++ {
		if len(game.Board[i]) != 3 {
			return domainmodel.Game{}, ErrInvalidBoardSize
		}
		for j := 0; j < 3; j++ {
			if game.Board[i][j] != int(domainmodel.Nought) && game.Board[i][j] != int(domainmodel.Cross) && game.Board[i][j] != int(domainmodel.Empty) {
				return domainmodel.Game{}, ErrInvalidCellValue
			}
			board[i][j] = domainmodel.Cell(game.Board[i][j])
		}
	}
	return domainmodel.Game{ID: game.ID, Board: board}, nil

}
func ToWeb(game domainmodel.Game) webmodel.Game {
	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = int(game.Board[i][j])
		}
	}
	return webmodel.Game{ID: game.ID, Board: board}
}
