package web

import (
	"errors"
	"strings"
	"tictactoe/domain"
	"tictactoe/web/webmodel"
)

var (
	ErrInvalidBoardSize = errors.New("board must be 3x3")
	ErrInvalidCellValue = errors.New("cell must be -1, 0 or 1")
	ErrInvalidLogin     = errors.New("invalid login input")
	ErrInvalidPassword  = errors.New("invalid password input")
)

func toDomain(game webmodel.Game) (domain.Game, error) {
	if len(game.Board) != 3 {
		return domain.Game{}, ErrInvalidBoardSize
	}
	var board domain.Board
	for i := 0; i < 3; i++ {
		if len(game.Board[i]) != 3 {
			return domain.Game{}, ErrInvalidBoardSize
		}
		for j := 0; j < 3; j++ {
			if game.Board[i][j] != int(domain.Nought) && game.Board[i][j] != int(domain.Cross) && game.Board[i][j] != int(domain.Empty) {
				return domain.Game{}, ErrInvalidCellValue
			}
			board[i][j] = domain.Cell(game.Board[i][j])
		}
	}
	return domain.Game{ID: game.ID, Board: board}, nil

}
func toWeb(game domain.Game) webmodel.Game {
	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = int(game.Board[i][j])
		}
	}
	return webmodel.Game{ID: game.ID, Board: board}
}

func signUpToDomain(request webmodel.SignUpRequest) (domain.SignUpRequest, error) {
	login := strings.TrimSpace(request.Login)
	if len(login) == 0 || strings.Contains(login, ":") {
		return domain.SignUpRequest{}, ErrInvalidLogin
	}
	if len(strings.TrimSpace(request.Password)) == 0 || len(request.Password) > 72 {
		return domain.SignUpRequest{}, ErrInvalidPassword
	}
	return domain.SignUpRequest{Login: login, Password: request.Password}, nil
}
