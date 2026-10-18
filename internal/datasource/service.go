package datasource

import (
	"errors"
	"tictactoe/domain"
)

type gameService struct {
	repo Repository
}

func NewService(repo Repository) domain.Service {
	return &gameService{repo: repo}

}

func (s *gameService) NextMove(game domain.Game) (domain.Game, error) {
	return domain.Game{}, nil

}

var (
	ErrBoardModified     = errors.New("previous moves were modified")
	ErrInvalidMovesCount = errors.New("exactly one new move is required")
)

func (s *gameService) ValidateBoard(game domain.Game) error {
	oldGame, err := s.repo.Get(game.ID)
	if err != nil {
		return err
	}
	newMoves := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if oldGame.Board[i][j] != domain.Empty && oldGame.Board[i][j] != game.Board[i][j] {
				return ErrBoardModified
			}

			if oldGame.Board[i][j] == domain.Empty && game.Board[i][j] != domain.Empty {
				if game.Board[i][j] != domain.Cross {
					return ErrBoardModified
				}
				newMoves++
			}
		}
	}
	if newMoves != 1 {
		return ErrInvalidMovesCount
	}
	return nil

}
func (s *gameService) IsGameOver(game domain.Game) (bool, domain.Cell) {
	winner := findWinner(game.Board)
	if winner != domain.Empty {
		return true, winner
	}
	return isBoardFull(game.Board), domain.Empty
}
func findWinner(board domain.Board) domain.Cell {
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][2] != domain.Empty {
			return board[i][0]
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[2][i] != domain.Empty {
			return board[0][i]
		}
	}
	if (board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != domain.Empty) ||
		(board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != domain.Empty) {
		return board[1][1]
	}
	return domain.Empty
}
func isBoardFull(board domain.Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				return false
			}
		}
	}
	return true

}
